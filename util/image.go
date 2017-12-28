package util

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	patternFromImage = "[fF][rR][oO][mM]\\s+([^/:]+)/?([^:]*):?(.*)"
	orgLaincloud     = "laincloud"
)

var (
	errFileNotInDockerfileContext = fmt.Errorf("file is not in Dockerfile context")
	regexFromImage                = regexp.MustCompile(patternFromImage)
)

// FromImage denotes `FROM` instruction in Dockerfile
type FromImage struct {
	Organization string
	Repository   string
	Tag          string
}

// Name return ${organization}/${repository}:${tag}
func (f FromImage) Name() string {
	return fmt.Sprintf("%s/%s:%s", f.Organization, f.Repository, f.Tag)
}

// Context is the path containing Dockerfile
func (f FromImage) Context() string {
	if f.Organization != orgLaincloud {
		return ""
	}

	return fmt.Sprintf("%s/%s", f.Repository, f.Tag)
}

// In check whether this image is in images
func (f FromImage) In(images map[string]Image) bool {
	_, ok := images[f.Name()]
	return ok
}

func readFromImage(dockerfilePath string) (*FromImage, error) {
	f, err := os.Open(dockerfilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errFileNotInDockerfileContext
		}

		return nil, err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	fromImageLine, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	return parseFromImage(fromImageLine)
}

// parseFromImage parse FromImage according to the FROM instruction in Dockerfile
func parseFromImage(fromImageLine string) (*FromImage, error) {
	xs := regexFromImage.FindStringSubmatch(fromImageLine)
	if len(xs) != 4 {
		return nil, fmt.Errorf("%s is not of pattern: %s", fromImageLine, patternFromImage)
	}

	organization, repository, tag := xs[1], xs[2], xs[3]
	if repository == "" {
		repository = organization
		organization = "library"
	}
	if tag == "" {
		tag = "latest"
	}
	return &FromImage{
		Organization: organization,
		Repository:   repository,
		Tag:          tag,
	}, nil
}

// Image contains necessary information of an image
type Image struct {
	Organization string
	Repository   string
	Tag          string
	From         FromImage
}

// Name return ${organization}/${repository}:${tag}
func (i Image) Name() string {
	return fmt.Sprintf("%s/%s:%s", i.Organization, i.Repository, i.Tag)
}

// Context is the path containing Dockerfile
func (i Image) Context() string {
	if i.Organization != orgLaincloud {
		return ""
	}

	return fmt.Sprintf("%s/%s", i.Repository, i.Tag)
}

// IsFrom check whether i is from image
func (i Image) IsFrom(image Image, allImages map[string]Image) bool {
	if i.From.Organization == image.Organization && i.From.Repository == image.Repository && i.From.Tag == image.Tag {
		return true
	}

	from, ok := allImages[i.From.Name()]
	if !ok {
		return false
	}

	return from.IsFrom(image, allImages)
}

// newImage create Image
func newImage(path string) (*Image, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}

	// Assume that symlink is `${repository}/${symTag} -> ${repository}/${tag}`,
	// then ${symTag} is an alias for ${tag}
	if info.Mode()&os.ModeSymlink != 0 {
		realPath, err := filepath.EvalSymlinks(path)
		if err != nil {
			return nil, err
		}

		fromImage, err := readFromImage(fmt.Sprintf("%s/Dockerfile", realPath))
		if err != nil {
			return nil, err
		}

		xs := strings.Split(path, "/")
		if len(xs) != 2 {
			return nil, fmt.Errorf("path: %s is not of pattern ${imageRepository}/${imageTag}", path)
		}

		return &Image{
			Organization: orgLaincloud,
			Repository:   xs[0],
			Tag:          xs[1],
			From:         *fromImage,
		}, nil
	}

	xs := strings.Split(path, "/")
	if len(xs) < 3 {
		return nil, errFileNotInDockerfileContext
	}

	fromImage, err := readFromImage(fmt.Sprintf("%s/%s/Dockerfile", xs[0], xs[1]))
	if err != nil {
		return nil, err
	}

	return &Image{
		Organization: orgLaincloud,
		Repository:   xs[0],
		Tag:          xs[1],
		From:         *fromImage,
	}, nil
}

// parseImages parse map[${imageName}]Image from files
func parseImages(files []string) (map[string]Image, error) {
	images := make(map[string]Image)
	for _, p := range files {
		if strings.HasPrefix(p, ".") { // ignore hide files
			continue
		}

		image, err := newImage(p)
		if err != nil {
			if os.IsNotExist(err) || err == errFileNotInDockerfileContext {
				continue
			}
			return nil, err
		}
		images[image.Name()] = *image
	}

	return images, nil
}
