package backend

import (
	"context"

"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type DatasetFolder struct {
Id     uuid.UUID
	Folder string
Cancel context.CancelFunc
}

// Select a folder using a native menu
func SelectFolder(context context.Context) (DatasetFolder, error) {
	dialogOptions := runtime.OpenDialogOptions{
	}

	folder, err := runtime.OpenDirectoryDialog(context, dialogOptions)
	if err != nil || folder == "" {
		return DatasetFolder{}, err
	}

	id := uuid.New()

	runtime.EventsEmit(context, "folder-added", id, folder)

	selected_folder := DatasetFolder{Folder: folder, Id: id}
	return selected_folder, nil
}
