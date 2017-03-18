package atlanticnet

type Image struct {
	Architecture string `json:"architecture"`
	DisplayName  string `json:"displayname"`
	Type         string `json:"image_type"`
	Id           string `json:"imageid"`
	OsType       string `json:"ostype"`
	Owner        string `json:"owner"`
	Platform     string `json:"platform"`
	Version      string `json:"version"`
}

type describeImageResponse struct {
	ImagesSet map[string]Image `json:"imagesset"`
	RequestId string           `json:"requestid"`
}

type describeImageDTO struct {
	BaseResponse
	Response describeImageResponse `json:"describe-imageresponse"`
}

func (c *client) DescribeImage(imageId string) ([]Image, error) {
	in := map[string]string{}
	if imageId != "" {
		in["imageid"] = imageId
	}
	out := describeImageDTO{}
	err := c.doRequest("describe-image", in, &out)
	if err != nil {
		return []Image{}, err
	} else if out.Error != nil {
		return []Image{}, out.Error
	}

	images := make([]Image, len(out.Response.ImagesSet))
	i := 0
	for _, image := range out.Response.ImagesSet {
		images[i] = image
		i += 1
	}

	return images, nil
}
