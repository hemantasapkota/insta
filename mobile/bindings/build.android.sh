rm -rf ../client/android/app/libs/bindings.aar
gomobile bind -v -target android
cp bindings.aar ../client/android/app/libs
