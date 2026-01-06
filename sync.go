package main

func SyncToClipboard(cb Clipboard, fileContent string) error {
	currentClipboard, err := cb.Read()
	if err != nil {
		return err
	}

	if currentClipboard != fileContent {
		return cb.Write(fileContent)
	}

	return nil
}
