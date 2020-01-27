package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/alphagov/travis-resource/common"
	. "github.com/alphagov/travis-resource/in/command"
	"github.com/alphagov/travis-resource/messager"
	"github.com/alphagov/travis-resource/model"
)

func main() {
	ctx := context.Background()
	mes := messager.GetMessager()
	if len(os.Args) <= 1 {
		mes.FatalIf("error in command argument", errors.New("you must pass a folder as a first argument"))
	}
	destinationFolder := os.Args[1]
	err := os.MkdirAll(destinationFolder, 0755)
	if err != nil {
		mes.FatalIf("creating destination", err)
	}
	var request model.InRequest
	err = json.NewDecoder(os.Stdin).Decode(&request)
	mes.FatalIf("failed to read request ", err)

	if request.Source.Repository == "" {
		mes.FatalIf("can't get build", errors.New("there is no repository set"))
	}

	travisClient, err := common.MakeTravisClient(ctx, request.Source)
	mes.FatalIf("failed to create travis client", err)

	inCommand := &InCommand{travisClient, request, destinationFolder, mes}
	build, err := inCommand.GetBuildInfo(ctx)
	mes.FatalIf("can't get build", err)

	err = inCommand.WriteInBuildInfoFile(build)
	err = inCommand.WriteInCommitRefFile(build)
	mes.FatalIf("can't create file build info", err)

	err = inCommand.DownloadLogs(ctx, build)
	mes.FatalIf("can't download logs", err)

	inCommand.SendResponse(build)
}
