package util

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const (
	patternFromImage = "\\s*[Ff][Rr][Oo][Mm]\\s+([^/:]+)/?([^:]*):?(.*)"
	patternTags      = "\\s*#\\s*[Tt][Aa][Gg][Ss]\\s+(.+)"
	orgLaincloud     = "laincloud"
)

var (
	errFileNotInDockerfileContext = fmt.Errorf("file is not in Dockerfile context")
	regexFromImage                = regexp.MustCompile(patternFromImage)
	regexTags                     = regexp.MustCompile(patternTags)
)

// FromImage denotes `FROM` instruction in Dockerfile
type FromImage struct {
	Organization string
	Repository   string
	Tag          string
}

// In check whether this image is in images
// images: map[${image.Context}]Image
func (f FromImage) In(images map[string]Image) bool {
	for _, image := range images {
		if f.Organization == image.Organization && f.Repository == image.Repository {
			for _, tag := range image.Tags {
				if f.Tag == tag {
					return true
				}
			}
		}
	}

	return false
}

// Context get the context of FromImage
// allImages: map[${image.Context}]Image
func (f FromImage) Context(allImages map[string]Image) string {
	for _, image := range allImages {
		if f.Organization == image.Organization && f.Repository == image.Repository {
			for _, tag := range image.Tags {
				if f.Tag == tag {
					return image.Context
				}
			}
		}
	}

	return ""
}

// parseFromImage parse FromImage according to the FROM instruction in Dockerfile
func parseFromImage(fromImageLine string) (*FromImage, error) {
	xs := regexFromImage.FindStringSubmatch(fromImageLine)
	if len(xs) != 4 {
		return nil, fmt.Errorf("%s is not of pattern %s", fromImageLine, patternFromImage)
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
	Context      string // The directory contains the Dockerfile
	From         FromImage
	Organization string
	Repository   string
	Tags         []string
}

// IsFrom check whether i is from image
// allImages: map[${image.Context}]Image
func (i Image) IsFrom(baseImage Image, allImages map[string]Image) bool {
	if i.From.Organization == baseImage.Organization && i.From.Repository == baseImage.Repository {
		for _, tag := range baseImage.Tags {
			if i.From.Tag == tag {
				return true
			}
		}
	}

	for _, image := range allImages {
		if i.From.Organization == image.Organization && i.From.Repository == image.Repository {
			for _, tag := range image.Tags {
				if i.From.Tag == tag {
					return image.IsFrom(baseImage, allImages)
				}
			}
		}
	}

	return false
}

// newImage create Image according to file
func newImage(file string) (*Image, error) {
	xs := strings.Split(file, "/")
	if len(xs) < 3 {
		return nil, errFileNotInDockerfileContext
	}

	f, err := os.Open(fmt.Sprintf("%s/%s/Dockerfile", xs[0], xs[1]))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errFileNotInDockerfileContext
		}

		return nil, err
	}
	defer f.Close()

	var fromImage *FromImage
	tags := []string{xs[1]}
	scanner := bufio.NewScanner(f)
	for i := 0; scanner.Scan(); i++ {
		if i == 0 && strings.HasPrefix(scanner.Text(), "#") {
			if tags, err = parseTags(scanner.Text()); err != nil {
				return nil, err
			}

			continue
		}

		if scanner.Text() == "" {
			continue
		}

		if fromImage, err = parseFromImage(scanner.Text()); err != nil {
			return nil, err
		}
		break
	}

	return &Image{
		Context:      fmt.Sprintf("%s/%s", xs[0], xs[1]),
		Organization: orgLaincloud,
		Repository:   xs[0],
		Tags:         tags,
		From:         *fromImage,
	}, nil
}

func parseTags(tagsLine string) ([]string, error) {
	xs := regexTags.FindStringSubmatch(tagsLine)
	if len(xs) != 2 {
		return nil, fmt.Errorf("%s is not of pattern %s", tagsLine, patternTags)
	}

	rawTags := strings.Split(strings.TrimSpace(xs[1]), " ")
	tags := make([]string, 0)
	for _, rawTag := range rawTags {
		if rawTag != "" {
			tags = append(tags, rawTag)
		}
	}
	return tags, nil
}

// getContext2Images get map[${image.Context}]Image from files
func GetContext2Images(files []string) (map[string]Image, error) {
	images := make(map[string]Image)
	for _, f := range files {
		if strings.HasPrefix(f, ".") { // ignore hide files
			continue
		}

		image, err := newImage(f)
		if err != nil {
			if err == errFileNotInDockerfileContext {
				continue
			}

			return nil, err
		}

		images[image.Context] = *image
	}

	return images, nil
}
