package main

import (
	"net/http"
	"testing"
)

func fakemarshal(v interface{}) ([]byte, error) {
    return []byte{}, errors.New("Marshalling failed")
}

func restoremarshal(replace func(v interface{}) ([]byte, error)) {
    jsonMarshal = replace
}

func TestPolicyDocumentToStr(t * testing.T){
    storedMarshal := jsonMarshal
    jsonMarshal = fakemarshal
    defer restoremarshal(storedMarshal)

    input := map[string]interface{} {
        "test": "test1",
    }
    tests := []struct {
        name string
        arg map[string]interface{}
        wantErr string
    }{
        {
            name: "Test if JSON Marshalling fails",
            arg: input,
            wantErr: "Marshalling failed",
        },
    }

    for _, tt := range tests {
        _, gotErr := policyDocumentToStr(tt.arg)
       if gotErr != nil && gotErr.Error() != tt.wantErr {
            t.Errorf("Expected %s but got %s", tt.wantErr, gotErr.Error())
       }
    }
}