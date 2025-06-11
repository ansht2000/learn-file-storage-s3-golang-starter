package main

import "testing"

func TestGetVideoAspectRatio(t *testing.T) {
	cases := []struct{
		filepath string
		expectedAspectRatio string
	}{
		{
			filepath: "./samples/boots-video-horizontal.mp4",
			expectedAspectRatio: "16:9",
		},
		{
			filepath: "./samples/boots-video-vertical.mp4",
			expectedAspectRatio: "9:16",
		},
	}

	for _, c := range cases {
		aspectRatio, _ := getVideoAspectRatio(c.filepath)
		if aspectRatio != c.expectedAspectRatio {
			t.Errorf(
				"expected aspect ratio: %s does not match returned aspect ratio: %s",
				c.expectedAspectRatio,
				aspectRatio,
			)
		}
	}
}