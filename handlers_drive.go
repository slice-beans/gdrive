package main

import (
	"gdrive/auth"
	"gdrive/cli"
	"gdrive/drive"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

const ClientId = "898818517276-lp8epk1bffo6aac3oruof6u7gido2rd8.apps.googleusercontent.com"
const ClientSecret = "oQrEZGMplNc3M7MVgHwgWv4H"
const TokenFilename = "token_v2.json"
const DefaultCacheFileName = "file_cache.json"

func listHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).List(drive.ListFilesArgs{
		Out:         os.Stdout,
		MaxFiles:    args.Int64("maxFiles"),
		NameWidth:   args.Int64("nameWidth"),
		Query:       args.String("query"),
		SortOrder:   args.String("sortOrder"),
		SkipHeader:  args.Bool("skipHeader"),
		SizeInBytes: args.Bool("sizeInBytes"),
		AbsPath:     args.Bool("absPath"),
	})
	checkErr(err)
}

func listChangesHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).ListChanges(drive.ListChangesArgs{
		Out:        os.Stdout,
		PageToken:  args.String("pageToken"),
		MaxChanges: args.Int64("maxChanges"),
		Now:        args.Bool("now"),
		NameWidth:  args.Int64("nameWidth"),
		SkipHeader: args.Bool("skipHeader"),
	})
	checkErr(err)
}

func downloadHandler(ctx cli.Context) {
	args := ctx.Args()
	checkDownloadArgs(args)
	err := newDrive(args).Download(drive.DownloadArgs{
		Out:       os.Stdout,
		Id:        args.String("fileId"),
		Force:     args.Bool("force"),
		Path:      args.String("path"),
		Delete:    args.Bool("delete"),
		Recursive: args.Bool("recursive"),
		Stdout:    args.Bool("stdout"),
		Progress:  progressWriter(args.Bool("noProgress")),
	})
	checkErr(err)
}

func downloadQueryHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).DownloadQuery(drive.DownloadQueryArgs{
		Out:       os.Stdout,
		Query:     args.String("query"),
		Force:     args.Bool("force"),
		Recursive: args.Bool("recursive"),
		Path:      args.String("path"),
		Progress:  progressWriter(args.Bool("noProgress")),
	})
	checkErr(err)
}

func downloadSyncHandler(ctx cli.Context) {
	args := ctx.Args()
	cachePath := filepath.Join(args.String("configDir"), DefaultCacheFileName)
	err := newDrive(args).DownloadSync(drive.DownloadSyncArgs{
		Out:              os.Stdout,
		Progress:         progressWriter(args.Bool("noProgress")),
		Path:             args.String("path"),
		RootId:           args.String("fileId"),
		DryRun:           args.Bool("dryRun"),
		DeleteExtraneous: args.Bool("deleteExtraneous"),
		Resolution:       conflictResolution(args),
		Comparer:         NewCachedMd5Comparer(cachePath),
	})
	checkErr(err)
}

func downloadRevisionHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).DownloadRevision(drive.DownloadRevisionArgs{
		Out:        os.Stdout,
		FileId:     args.String("fileId"),
		RevisionId: args.String("revId"),
		Force:      args.Bool("force"),
		Stdout:     args.Bool("stdout"),
		Path:       args.String("path"),
		Progress:   progressWriter(args.Bool("noProgress")),
	})
	checkErr(err)
}

func uploadHandler(ctx cli.Context) {
	args := ctx.Args()
	checkUploadArgs(args)
	err := newDrive(args).Upload(drive.UploadArgs{
		Out:       os.Stdout,
		Progress:  progressWriter(args.Bool("noProgress")),
		Path:      args.String("path"),
		Name:      args.String("name"),
		Parents:   args.StringSlice("parent"),
		Mime:      args.String("mime"),
		Recursive: args.Bool("recursive"),
		Share:     args.Bool("share"),
		Delete:    args.Bool("delete"),
		ChunkSize: args.Int64("chunksize"),
	})
	checkErr(err)
}

func uploadStdinHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).UploadStream(drive.UploadStreamArgs{
		Out:       os.Stdout,
		In:        os.Stdin,
		Name:      args.String("name"),
		Parents:   args.StringSlice("parent"),
		Mime:      args.String("mime"),
		Share:     args.Bool("share"),
		ChunkSize: args.Int64("chunksize"),
		Progress:  progressWriter(args.Bool("noProgress")),
	})
	checkErr(err)
}

func uploadSyncHandler(ctx cli.Context) {
	args := ctx.Args()
	cachePath := filepath.Join(args.String("configDir"), DefaultCacheFileName)
	err := newDrive(args).UploadSync(drive.UploadSyncArgs{
		Out:              os.Stdout,
		Progress:         progressWriter(args.Bool("noProgress")),
		Path:             args.String("path"),
		RootId:           args.String("fileId"),
		DryRun:           args.Bool("dryRun"),
		DeleteExtraneous: args.Bool("deleteExtraneous"),
		ChunkSize:        args.Int64("chunksize"),
		Resolution:       conflictResolution(args),
		Comparer:         NewCachedMd5Comparer(cachePath),
	})
	checkErr(err)
}

func updateHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).Update(drive.UpdateArgs{
		Out:       os.Stdout,
		Id:        args.String("fileId"),
		Path:      args.String("path"),
		Name:      args.String("name"),
		Parents:   args.StringSlice("parent"),
		Mime:      args.String("mime"),
		Progress:  progressWriter(args.Bool("noProgress")),
		ChunkSize: args.Int64("chunksize"),
	})
	checkErr(err)
}

func infoHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).Info(drive.FileInfoArgs{
		Out:         os.Stdout,
		Id:          args.String("fileId"),
		SizeInBytes: args.Bool("sizeInBytes"),
	})
	checkErr(err)
}

func importHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).Import(drive.ImportArgs{
		Out:      os.Stdout,
		Path:     args.String("path"),
		Parents:  args.StringSlice("parent"),
		Progress: progressWriter(args.Bool("noProgress")),
	})
	checkErr(err)
}

func exportHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).Export(drive.ExportArgs{
		Out:        os.Stdout,
		Id:         args.String("fileId"),
		Mime:       args.String("mime"),
		PrintMimes: args.Bool("printMimes"),
		Force:      args.Bool("force"),
	})
	checkErr(err)
}

func listRevisionsHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).ListRevisions(drive.ListRevisionsArgs{
		Out:         os.Stdout,
		Id:          args.String("fileId"),
		NameWidth:   args.Int64("nameWidth"),
		SizeInBytes: args.Bool("sizeInBytes"),
		SkipHeader:  args.Bool("skipHeader"),
	})
	checkErr(err)
}

func mkdirHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).Mkdir(drive.MkdirArgs{
		Out:     os.Stdout,
		Name:    args.String("name"),
		Parents: args.StringSlice("parent"),
	})
	checkErr(err)
}

func shareHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).Share(drive.ShareArgs{
		Out:          os.Stdout,
		FileId:       args.String("fileId"),
		Role:         args.String("role"),
		Type:         args.String("type"),
		Email:        args.String("email"),
		Discoverable: args.Bool("discoverable"),
	})
	checkErr(err)
}

func shareListHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).ListPermissions(drive.ListPermissionsArgs{
		Out:    os.Stdout,
		FileId: args.String("fileId"),
	})
	checkErr(err)
}

func shareRevokeHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).RevokePermission(drive.RevokePermissionArgs{
		Out:          os.Stdout,
		FileId:       args.String("fileId"),
		PermissionId: args.String("permissionId"),
	})
	checkErr(err)
}

func deleteHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).Delete(drive.DeleteArgs{
		Out:       os.Stdout,
		Id:        args.String("fileId"),
		Recursive: args.Bool("recursive"),
	})
	checkErr(err)
}

func listSyncHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).ListSync(drive.ListSyncArgs{
		Out:        os.Stdout,
		SkipHeader: args.Bool("skipHeader"),
	})
	checkErr(err)
}

func listRecursiveSyncHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).ListRecursiveSync(drive.ListRecursiveSyncArgs{
		Out:         os.Stdout,
		RootId:      args.String("fileId"),
		SkipHeader:  args.Bool("skipHeader"),
		PathWidth:   args.Int64("pathWidth"),
		SizeInBytes: args.Bool("sizeInBytes"),
		SortOrder:   args.String("sortOrder"),
	})
	checkErr(err)
}

func deleteRevisionHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).DeleteRevision(drive.DeleteRevisionArgs{
		Out:        os.Stdout,
		FileId:     args.String("fileId"),
		RevisionId: args.String("revId"),
	})
	checkErr(err)
}

func aboutHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).About(drive.AboutArgs{
		Out:         os.Stdout,
		SizeInBytes: args.Bool("sizeInBytes"),
	})
	checkErr(err)
}

func aboutImportHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).AboutImport(drive.AboutImportArgs{
		Out: os.Stdout,
	})
	checkErr(err)
}

func aboutExportHandler(ctx cli.Context) {
	args := ctx.Args()
	err := newDrive(args).AboutExport(drive.AboutExportArgs{
		Out: os.Stdout,
	})
	checkErr(err)
}

func getOauthClient(args cli.Arguments) (*http.Client, error) {
	if args.String("refreshToken") != "" && args.String("accessToken") != "" {
		ExitF("Access token not needed when refresh token is provided")
	}

	if args.String("refreshToken") != "" {
		return auth.NewRefreshTokenClient(ClientId, ClientSecret, args.String("refreshToken")), nil
	}

	if args.String("accessToken") != "" {
		return auth.NewAccessTokenClient(ClientId, ClientSecret, args.String("accessToken")), nil
	}

	configDir := getConfigDir(args)
	tokenPath := ConfigFilePath(configDir, TokenFilename)
	return auth.NewFileSourceClient(ClientId, ClientSecret, tokenPath, authCodePrompt)
}

func getConfigDir(args cli.Arguments) string {
	// Use dir from environment var if present
	if os.Getenv("GDRIVE_CONFIG_DIR") != "" {
		return os.Getenv("GDRIVE_CONFIG_DIR")
	}
	return args.String("configDir")
}

func newDrive(args cli.Arguments) *drive.Drive {
	oauth, err := getOauthClient(args)
	if err != nil {
		ExitF("Failed getting oauth client: %s", err.Error())
	}

	client, err := drive.New(oauth)
	if err != nil {
		ExitF("Failed getting drive: %s", err.Error())
	}

	return client
}

func authCodePrompt(url string) func() string {
	return func() string {
		fmt.Println("Authentication needed")
		fmt.Println("Go to the following url in your browser:")
		fmt.Printf("%s\n\n", url)
		fmt.Print("Enter verification code: ")

		var code string
		if _, err := fmt.Scan(&code); err != nil {
			fmt.Printf("Failed reading code: %s", err.Error())
		}
		return code
	}
}

func progressWriter(discard bool) io.Writer {
	if discard {
		return ioutil.Discard
	}
	return os.Stderr
}

func conflictResolution(args cli.Arguments) drive.ConflictResolution {
	keepLocal := args.Bool("keepLocal")
	keepRemote := args.Bool("keepRemote")
	keepLargest := args.Bool("keepLargest")

	if (keepLocal && keepRemote) || (keepLocal && keepLargest) || (keepRemote && keepLargest) {
		ExitF("Only one conflict resolution flag can be given")
	}

	if keepLocal {
		return drive.KeepLocal
	}

	if keepRemote {
		return drive.KeepRemote
	}

	if keepLargest {
		return drive.KeepLargest
	}

	return drive.NoResolution
}

func checkUploadArgs(args cli.Arguments) {
	if args.Bool("recursive") && args.Bool("delete") {
		ExitF("--delete is not allowed for recursive uploads")
	}

	if args.Bool("recursive") && args.Bool("share") {
		ExitF("--share is not allowed for recursive uploads")
	}
}

func checkDownloadArgs(args cli.Arguments) {
	if args.Bool("recursive") && args.Bool("delete") {
		ExitF("--delete is not allowed for recursive downloads")
	}
}
