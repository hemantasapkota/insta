rm -rf ../client/ios/Bindings.framework
gomobile bind -v -target ios
mv Bindings.framework ../client/ios/
