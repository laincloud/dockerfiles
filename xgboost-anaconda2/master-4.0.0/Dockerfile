FROM laincloud/debian:stretch
LABEL maintainer="Ren Jingsi <jingsiren2@creditease.cn>"

ENV XGBOOST_VERSION=0.60 ANACONDA2_VERSION=4.0.0 \
    PATH=/opt/conda/bin:$PATH LANG=C.UTF-8
RUN apt-get update && apt-get install -y git gcc build-essential \
    wget bzip2 ca-certificates libglib2.0-0 libxext6 libsm6 \
    libxrender1 mercurial subversion

RUN cd /opt; git clone --recursive https://github.com/dmlc/xgboost && \
    cd /opt/xgboost; make -j4

RUN cd /; echo 'export PATH=/opt/conda/bin:$PATH' > /etc/profile.d/conda.sh && \
    wget https://repo.continuum.io/archive/Anaconda2-$ANACONDA2_VERSION-Linux-x86_64.sh && \
    /bin/bash /Anaconda2-$ANACONDA2_VERSION-Linux-x86_64.sh -b -p /opt/conda && \
    rm /Anaconda2-$ANACONDA2_VERSION-Linux-x86_64.sh

RUN conda install -y -c https://conda.anaconda.org/anaconda setuptools && \
    cd /opt/xgboost/python-package/ && python setup.py install && pip install gunicorn

# Add a notebook profile.
RUN mkdir -p -m 700 /root/.jupyter/ && \
    echo "c.NotebookApp.ip = '*'" >> /root/.jupyter/jupyter_notebook_config.py

VOLUME /notebooks
WORKDIR /notebooks

EXPOSE 8888

ENTRYPOINT ["/app-init"]
CMD ["jupyter", "notebook"]
