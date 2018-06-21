#!/bin/bash
#
# An entrypoint script for Artifactory to allow custom setup before server starts
#

: ${ARTIFACTORY_USER_NAME:=artifactory}
: ${ARTIFACTORY_USER_ID:=1030}
: ${ARTIFACTORY_HOME:=/opt/jfrog/artifactory}
: ${ARTIFACTORY_DATA:=/var/opt/jfrog/artifactory}
: ${ACCESS_ETC_FOLDER=${ARTIFACTORY_DATA}/access/etc}
ART_ETC=$ARTIFACTORY_DATA/etc

: ${RECOMMENDED_MAX_OPEN_FILES:=32000}
: ${MIN_MAX_OPEN_FILES:=10000}

: ${RECOMMENDED_MAX_OPEN_PROCESSES:=1024}

export ARTIFACTORY_PID=${ARTIFACTORY_HOME}/run/artifactory.pid

logger() {
    DATE_TIME=$(date +"%Y-%m-%d %H:%M:%S")
    if [ -z "$CONTEXT" ]
    then
        CONTEXT=$(caller)
    fi
    MESSAGE=$1
    CONTEXT_LINE=$(echo "$CONTEXT" | awk '{print $1}')
    CONTEXT_FILE=$(echo "$CONTEXT" | awk -F"/" '{print $NF}')
    printf "%s %05s %s %s\n" "$DATE_TIME" "[$CONTEXT_LINE" "$CONTEXT_FILE]" "$MESSAGE"
    CONTEXT=
}

errorExit () {
    logger "ERROR: $1"; echo
    exit 1
}

warn () {
    logger "WARNING: $1"
}

# Print on container startup information about Dockerfile location
printDockerFileLocation() {
    logger "Dockerfile for this image can found inside the container."
    logger "To view the Dockerfile: 'cat /docker/artifactory-oss/Dockerfile.artifactory'."
}

# Check the max open files and open processes set on the system
checkULimits () {
    logger "Checking open files and processes limits"

    CURRENT_MAX_OPEN_FILES=$(ulimit -n)
    logger "Current max open files is $CURRENT_MAX_OPEN_FILES"

    if [ ${CURRENT_MAX_OPEN_FILES} != "unlimited" ] && [ "$CURRENT_MAX_OPEN_FILES" -lt "$RECOMMENDED_MAX_OPEN_FILES" ]; then
        if [ "$CURRENT_MAX_OPEN_FILES" -lt "$MIN_MAX_OPEN_FILES" ]; then
            errorExit "Max number of open files $CURRENT_MAX_OPEN_FILES, is too low. Cannot run Artifactory!"
        fi

        warn "Max number of open files $CURRENT_MAX_OPEN_FILES is low!"
        warn "You should add the parameter '--ulimit nofile=${RECOMMENDED_MAX_OPEN_FILES}:${RECOMMENDED_MAX_OPEN_FILES}' to your the 'docker run' command."
    fi

    CURRENT_MAX_OPEN_PROCESSES=$(ulimit -u)
    logger "Current max open processes is $CURRENT_MAX_OPEN_PROCESSES"

    if [ "$CURRENT_MAX_OPEN_PROCESSES" != "unlimited" ] && [ "$CURRENT_MAX_OPEN_PROCESSES" -lt "$RECOMMENDED_MAX_OPEN_PROCESSES" ]; then
        warn "Max number of processes $CURRENT_MAX_OPEN_PROCESSES is too low!"
        warn "You should add the parameter '--ulimit noproc=${RECOMMENDED_MAX_OPEN_PROCESSES}:${RECOMMENDED_MAX_OPEN_PROCESSES}' to your the 'docker run' command."
    fi
}

# Check that data dir is mounted and warn if not
checkMounts () {
    logger "Checking if $ARTIFACTORY_DATA is mounted"
    mount | grep ${ARTIFACTORY_DATA} > /dev/null
    if [ $? -ne 0 ]; then
        warn "Artifactory data directory ($ARTIFACTORY_DATA) is not mounted from the host. This means that all data and configurations will be lost once container is removed!"
    else
        logger "$ARTIFACTORY_DATA is mounted"
    fi
}

# In case data dirs are missing or not mounted, need to create them
setupDataDirs () {
    logger "Setting up data directories if missing"
    [ -d ${ARTIFACTORY_DATA}/etc ]    || mkdir -p ${ARTIFACTORY_DATA}/etc    || errorExit "Creating ${ARTIFACTORY_DATA}/etc folder failed"
    [ -d ${ARTIFACTORY_DATA}/data ]   || mkdir -p ${ARTIFACTORY_DATA}/data   || errorExit "Creating ${ARTIFACTORY_DATA}/data folder failed"
    [ -d ${ARTIFACTORY_DATA}/logs ]   || mkdir -p ${ARTIFACTORY_DATA}/logs   || errorExit "Creating ${ARTIFACTORY_DATA}/logs folder failed"
    [ -d ${ARTIFACTORY_DATA}/backup ] || mkdir -p ${ARTIFACTORY_DATA}/backup || errorExit "Creating ${ARTIFACTORY_DATA}/backup folder failed"
    [ -d ${ARTIFACTORY_DATA}/access ] || mkdir -p ${ARTIFACTORY_DATA}/access || errorExit "Creating ${ARTIFACTORY_DATA}/access folder failed"
    [ -d ${ARTIFACTORY_HOME}/run ]    || mkdir -p ${ARTIFACTORY_HOME}/run    || errorExit "Creating ${ARTIFACTORY_HOME}/run folder failed"
}

# Create the Artifactory user (support passing name and id as parameters)
setupArtUser () {
    logger "Create $ARTIFACTORY_USER_NAME user if missing"
    id -u ${ARTIFACTORY_USER_NAME} > /dev/null 2>&1
    if [ $? -ne 0 ]; then
        logger "User does not exist. Creating it..."
        useradd -M -s /usr/sbin/nologin --uid ${ARTIFACTORY_USER_ID} --user-group ${ARTIFACTORY_USER_NAME} || errorExit "Creating user ${ARTIFACTORY_USER_NAME} failed"
    else
        logger "User ${ARTIFACTORY_USER_NAME} already exists"
    fi
}

setAccessCreds() {
    ACCESS_SOURCE_IP_ALLOWED=${ACCESS_SOURCE_IP_ALLOWED:-127.0.0.1}
    ACCESS_CREDS_FILE=${ACCESS_ETC_FOLDER}/bootstrap.creds
    if [ ! -z "${ACCESS_USER}" ] && [ ! -z "${ACCESS_PASSWORD}" ] && [ ! -f ${ACCESS_CREDS_FILE} ] ; then
        logger "Creating bootstrap.creds using ACCESS_USER and ACCESS_PASSWORD env variables"
        mkdir -p ${ACCESS_ETC_FOLDER} || errorExit "Couldn't create ${ACCESS_ETC_FOLDER}"
        echo "${ACCESS_USER}@${ACCESS_SOURCE_IP_ALLOWED}=${ACCESS_PASSWORD}" > ${ACCESS_CREDS_FILE}
        chmod 600 ${ACCESS_CREDS_FILE} || errorExit "Setting permission on ${ACCESS_CREDS_FILE} failed"
        chown ${ARTIFACTORY_USER_NAME}:${ARTIFACTORY_USER_NAME} ${ACCESS_CREDS_FILE} || errorExit "Setting ownership on ${ACCESS_CREDS_FILE} failed"
    fi
}

setMasterKey() {
    ARTIFACTORY_SECURITY_FOLDER=${ART_ETC}/security
    ARTIFACTORY_MASTER_KEY_FILE=${ARTIFACTORY_SECURITY_FOLDER}/master.key
    if [ ! -z "${ARTIFACTORY_MASTER_KEY}" ] ; then
        logger "Creating master.key using ARTIFACTORY_MASTER_KEY environment variable"
        mkdir -p ${ARTIFACTORY_SECURITY_FOLDER} || errorExit "Creating folder ${ARTIFACTORY_SECURITY_FOLDER} failed"
        echo "${ARTIFACTORY_MASTER_KEY}" > "${ARTIFACTORY_MASTER_KEY_FILE}"
    fi

    if [ -f "${ARTIFACTORY_MASTER_KEY_FILE}" ] ; then
        chmod 600 ${ARTIFACTORY_MASTER_KEY_FILE} || errorExit "Setting permission on ${ARTIFACTORY_MASTER_KEY_FILE} failed"
        chown ${ARTIFACTORY_USER_NAME}:${ARTIFACTORY_USER_NAME} ${ARTIFACTORY_MASTER_KEY_FILE} || errorExit "Setting ownership on ${ARTIFACTORY_MASTER_KEY_FILE} failed"
    fi
}

# Check and set permissions on ARTIFACTORY_HOME and ARTIFACTORY_DATA
setupPermissions () {
    # ARTIFACTORY_HOME folder
    logger "Checking permissions on ${ARTIFACTORY_HOME}"
    STAT=( $(stat -Lc "%U %G" ${ARTIFACTORY_HOME}) )
    USER=${STAT[0]}
    GROUP=${STAT[1]}

    if [[ ${USER} != "$ARTIFACTORY_USER_NAME" ]] || [[ ${GROUP} != "$ARTIFACTORY_USER_NAME"  ]] ; then
        logger "$ARTIFACTORY_HOME is owned by $USER:$GROUP. Setting to $ARTIFACTORY_USER_NAME:$ARTIFACTORY_USER_NAME."
        chown -R ${ARTIFACTORY_USER_NAME}:${ARTIFACTORY_USER_NAME} ${ARTIFACTORY_HOME} || errorExit "Setting ownership on $ARTIFACTORY_HOME failed"
    else
        logger "$ARTIFACTORY_HOME is already owned by $ARTIFACTORY_USER_NAME:$ARTIFACTORY_USER_NAME."
    fi

    # ARTIFACTORY_DATA folder
    logger "Checking permissions on $ARTIFACTORY_DATA"
    STAT=( $(stat -Lc "%U %G" ${ARTIFACTORY_DATA}) )
    USER=${STAT[0]}
    GROUP=${STAT[1]}

    if [[ ${USER} != "$ARTIFACTORY_USER_NAME" ]] || [[ ${GROUP} != "$ARTIFACTORY_USER_NAME"  ]] ; then
        logger "$ARTIFACTORY_DATA is owned by $USER:$GROUP. Setting to $ARTIFACTORY_USER_NAME:$ARTIFACTORY_USER_NAME."
        chown -R ${ARTIFACTORY_USER_NAME}:${ARTIFACTORY_USER_NAME} ${ARTIFACTORY_DATA} || errorExit "Setting ownership on $ARTIFACTORY_DATA failed"
    else
        logger "$ARTIFACTORY_DATA is already owned by $ARTIFACTORY_USER_NAME:$ARTIFACTORY_USER_NAME."
    fi
}

# Wait for DB port to be accessible
waitForDB () {
    local PROPS_FILE=$1
    local DB_TYPE=$2

    [ -f "$PROPS_FILE" ] || errorExit "$PROPS_FILE does not exist"

    local DB_HOST_PORT=
    local TIMEOUT=30
    local COUNTER=0

    # Extract DB host and port
    case "$DB_TYPE" in
        postgresql|mysql)
            DB_HOST_PORT=$(grep -e '^url=' "$PROPS_FILE" | sed -e 's,^.*:\/\/\(.*\)\/.*,\1,g' | tr ':' '/')
        ;;
        oracle)
            DB_HOST_PORT=$(grep -e '^url=' "$PROPS_FILE" | sed -e 's,.*@\(.*\):.*,\1,g' | tr ':' '/')
        ;;
        mssql)
            DB_HOST_PORT=$(grep -e '^url=' "$PROPS_FILE" | sed -e 's,^.*:\/\/\(.*\);databaseName.*,\1,g' | tr ':' '/')
        ;;
        *)
            errorExit "DB_TYPE $DB_TYPE not supported"
        ;;
    esac

    logger "Waiting for DB $DB_TYPE to be ready on $DB_HOST_PORT within $TIMEOUT seconds"

    while [ $COUNTER -lt $TIMEOUT ]; do
        (</dev/tcp/$DB_HOST_PORT) 2>/dev/null
        if [ $? -eq 0 ]; then
            logger "DB $DB_TYPE up in $COUNTER seconds"
            return 1
        else
            logger "."
            sleep 1
        fi
        let COUNTER=$COUNTER+1
    done

    return 0
}

# Check DB type configurations before starting Artifactory
setDBConf () {
       # Set DB_HOST
    if [ -z "$DB_HOST" ]; then
        DB_HOST=$DB_TYPE
    fi
    logger "DB_HOST is set to $DB_HOST"

    logger "Checking if need to copy $DB_TYPE configuration"
    # If already exists, just make sure it's configured for postgres
    if [ -f ${DB_PROPS} ]; then
        logger "${DB_PROPS} already exists. Making sure it's set to $DB_TYPE... "
        grep type=$DB_TYPE ${DB_PROPS} > /dev/null
        if [ $? -eq 0 ]; then
            logger "${DB_PROPS} already set to $DB_TYPE"
        else
            errorExit "${DB_PROPS} already exists and is set to a DB different than $DB_TYPE"
        fi
    else
        NEED_COPY=true
    fi

    # On a new install and startup, need to make the initial copy before Artifactory starts
    if [ "$NEED_COPY" == "true" ]; then
        logger "Copying $DB_TYPE configuration... "
        cp ${ARTIFACTORY_HOME}/misc/db/$DB_TYPE.properties ${DB_PROPS} || errorExit "Copying $ARTIFACTORY_HOME/misc/db/$DB_TYPE.properties to ${DB_PROPS} failed"
        chown ${ARTIFACTORY_USER_NAME}: ${DB_PROPS} || errorExit "Change owner of ${DB_PROPS} to ${ARTIFACTORY_USER_NAME} failed"

        sed -i "s/localhost/$DB_HOST/g" ${DB_PROPS}

        # Set custom DB parameters if specified
        if [ ! -z "$DB_URL" ]; then
            logger "Setting DB_URL to $DB_URL"
            sed -i "s|url=.*|url=$DB_URL|g" ${DB_PROPS}
        fi
        if [ ! -z "$DB_USER" ]; then
            logger "Setting DB_USER to $DB_USER"
            sed -i "s/username=.*/username=$DB_USER/g" ${DB_PROPS}
        fi
        if [ ! -z "$DB_PASSWORD" ]; then
            logger "Setting DB_PASSWORD to **********"
            sed -i "s/password=.*/password=$DB_PASSWORD/g" ${DB_PROPS}
        fi
        if [ ! -z "$DB_PORT" ]; then
            logger "Setting DB_PORT to $DB_PORT"
            case "$DB_TYPE" in
            mysql|postgresql)
                oldPort=$(grep -E "(url).*" ${DB_PROPS}  | awk -F":" '{print $4}' | awk -F"/" '{print $1}')
            ;;
            oracle)
                oldPort=$(grep -E "(url).*" ${DB_PROPS} | awk -F":" '{print $5}')
            ;;
            mssql)
                oldPort=$(grep -E "(url).*" ${DB_PROPS}  | awk -F":" '{print $4}' | awk -F";" '{print $1}')
            ;;
            esac
               sed -i "s/$oldPort/$DB_PORT/g" ${DB_PROPS}
        fi
        if [ ! -z "$DB_HOST" ]; then
            logger "Setting DB_HOST to $DB_HOST"
            case "$DB_TYPE" in
            mysql|postgresql|mssql)
                oldHost=$(grep -E "(url).*" ${DB_PROPS} | awk -F"//" '{print $2}' | awk -F":" '{print $1}')
            ;;
            oracle)
                oldHost=$(grep -E "(url).*" ${DB_PROPS} | awk -F"@" '{print $2}' | awk -F":" '{print $1}')
            ;;
            esac
            sed -i "s/$oldHost/$DB_HOST/g" ${DB_PROPS}
        fi
    fi
}

# Set and configure DB type
setDBType () {
    logger "Checking DB_TYPE"

    if [ ! -z "$DB_TYPE" ]; then
        logger "DB_TYPE is set to $DB_TYPE"
        NEED_COPY=false
        DB_PROPS=${ART_ETC}/db.properties

        case "$DB_TYPE" in
            postgresql)
                if ! ls $ARTIFACTORY_HOME/tomcat/lib/postgresql-*.jar 1> /dev/null 2>&1; then
                    errorExit "No postgresql connector found"
                fi
                setDBConf
            ;;
            mysql)
                if ! ls $ARTIFACTORY_HOME/tomcat/lib/mysql-connector-java*.jar 1> /dev/null 2>&1; then
                    errorExit "No mysql connector found"
                fi
                setDBConf
            ;;
            oracle)
                if ! ls $ARTIFACTORY_HOME/tomcat/lib/ojdb*.jar 1> /dev/null 2>&1; then
                    errorExit "No oracle ojdbc driver found"
                fi
                setDBConf
            ;;
            mssql)
                if ! ls $ARTIFACTORY_HOME/tomcat/lib/sqljdbc*.jar 1> /dev/null 2>&1; then
                    errorExit "No mssql connector found"
                fi
                setDBConf
            ;;
            *)
                errorExit "DB_TYPE $DB_TYPE not supported"
            ;;
        esac

        # Wait for DB
        # On slow systems, when working with docker-compose, the DB container might be up,
        # but not ready to accept connections when Artifactory is already trying to access it.
        if [[ ! "$HA_IS_PRIMARY" =~ false ]]; then
            waitForDB "$DB_PROPS" "$DB_TYPE"
            [ $? -eq 1 ] || errorExit "DB $DB_TYPE failed to start in the given time"
        fi
    else
        logger "DB_TYPE not set. Artifactory will use built in Derby DB"
    fi
}

addExtraJavaArgs() {
    if [ ! -z "$EXTRA_JAVA_OPTIONS" ] && [ ! -f ${ARTIFACTORY_HOME}/bin/artifactory.default.origin ] ; then
        cp -v ${ARTIFACTORY_HOME}/bin/artifactory.default ${ARTIFACTORY_HOME}/bin/artifactory.default.origin
        echo "export JAVA_OPTIONS=\"\$JAVA_OPTIONS $EXTRA_JAVA_OPTIONS\"" >> ${ARTIFACTORY_HOME}/bin/artifactory.default
    fi
}

terminate () {
    echo -e "\nTerminating Artifactory"
    ${ARTIFACTORY_HOME}/bin/artifactory.sh stop
}

# Catch Ctrl+C and other termination signals to try graceful shutdown
trap terminate SIGINT SIGTERM SIGHUP

echo -e "\nPreparing to run Artifactory in Docker"
echo "====================================="

printDockerFileLocation
checkULimits
checkMounts
setupDataDirs
setupArtUser
setAccessCreds
setMasterKey
setupPermissions
setDBType
addExtraJavaArgs

echo -e "\n=====================================\n"

# Run Artifactory as ARTIFACTORY_USER_NAME user
exec gosu ${ARTIFACTORY_USER_NAME} ${ARTIFACTORY_HOME}/bin/artifactory.sh &
art_pid=$!

echo ${art_pid} > ${ARTIFACTORY_PID}

wait ${art_pid}
