# Insta - Power tool for Instagram.

Insta is a command line power tool for Instagram. With it, you can like, follow, comment and query Instagram.

# Motivation

There're plenty of Instagram automation tools floating around the web and app stores. Most of them ask for your username and password which is a big safety risk.

With insta, you can automate engagement yourself.

# How do i install it ?

* ``` git clone https://github.com/hemantasapkota/insta ```
* ``` go  get . && go build ```

# Usage

You can specify credentials inline or via a file.

* Inline credentials
``` insta -username user1  -password pass1 ```

* Create a **.credentials.yaml** file. Sample content below.

```yaml
acc1:
    username: user1 
    password: pass1
    
acc2:
    username: user2
    password: pass2

```

``` insta -account acc1 ```

# Show me a demo
[![asciicast](https://asciinema.org/a/8lfvvdi1gm6nd1t0rydfcmkci.png)](https://asciinema.org/a/8lfvvdi1gm6nd1t0rydfcmkci)

# How does it work ?

    #TODO

# Features

* Interactive console
* Manage multiple accounts
* Database log

# Roadmap

* Console app
* Mobile apps ( IOS & Android )
* Chrome extension

# Contribution Guidelines

   TODO

