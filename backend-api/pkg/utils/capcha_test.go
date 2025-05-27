package utils

import (
	"fmt"
	"testing"
	"time"
	// "strings"
)

func TestCaptcha(t *testing.T) {
	var (
		id, b64s, answer string
		err              error
		captcha          *Captcha
	)

	captcha = &Captcha{
		Enabled: true,
		Clear:   false,

		Height:     96,
		Width:      160,
		NoiseCount: 2,
		Length:     5,
		Source:     "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",

		CollectNum: 1000,
		Duration:   3 * time.Minute,
	}

	if err = captcha.Setup(); err != nil {
		t.Fatal(err)
	}

	if id, b64s, answer, err = captcha.Generate(); err != nil {
		t.Fatal(err)
	}

	fmt.Printf("==> 1. id=%s, answer=%s\n", id, answer)
	fmt.Printf("==> 2. base64(%dk)=%s...\n", len(b64s)/1024, b64s[:100])

	if !captcha.Verify(id, answer, true) {
		t.Fatal("not match")
	}

	if captcha.Verify("", answer, true) {
		t.Fatal("matched empty")
	}
}
