package main

import (
	"fmt"
	faker "github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	for n := 0; n < 120; n++ {
		t.Run(fmt.Sprintf("TestEncryptDecrypt-%v", n), func(t *testing.T) {
			t.Parallel()
			key := faker.Password()
			msg := faker.Paragraph() + faker.Email() + faker.TollFreePhoneNumber()

			for y := 0; y < 10; y++ {
				msg = msg + "(&**(#@;[;s[da;dsa"
				msg = msg + faker.Paragraph()
				msg = msg + "\t \n aaaa  asfmaskcBAHDB*!G*&HQ   (*"
			}

			b, err := encrypt([]byte(msg), key)
			if err != nil {
				t.Fatal(err)
			}

			p, err := decrypt(b, key)
			if err != nil {
				t.Fatal(err)
			}

			t.Logf("Total character %v\n", len(msg))
			assert.Equal(t, string(p), msg)
		})
	}
}
