set -e

# install packages
yarn install

# instal go deps
go install github.com/GeertJohan/go.rice
go install github.com/GeertJohan/go.rice/rice

# Builds CSS from SCSS
echo "gulp Gulp GULP"
gulp sass

# Removes mac shitty things
find assets/ -type f -name '.DS_Store' -delete

# Put assets into the binary
rice embed-go -i ./managers/assetsmgr

# Clean up data so it passes linter
gofmt -s -w ./managers/assetsmgr/rice-box.go 

echo "THAT WAS EASY!"