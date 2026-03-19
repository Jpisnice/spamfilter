# SpamFilter (Go)

SpamFilter is a small Go hobby project that experiments with spam/ham email classification using a simple bag-of-words approach on an Enron-style dataset; it is built for learning and testing ideas, not for real-world or production use.

To download the dataset zip directly into this project root as `enron-spam.zip`, run: `curl -L -o enron-spam.zip https://www.kaggle.com/api/v1/datasets/download/purusinghvi/email-spam-classification-dataset` (requires a configured Kaggle API token/account).

After download, unzip into `enron-spam/` in the same project root using either `unzip enron-spam.zip -d enron-spam` (bash/macOS/Linux) or `Expand-Archive -Path .\enron-spam.zip -DestinationPath .\enron-spam -Force` (PowerShell), then run the app with `go run main.go` (or `go run ./...`).
