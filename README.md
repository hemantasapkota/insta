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
[![asciicast](https://asciinema.org/a/6ksjrq7404k4877thbv5ega9r.png)](https://asciinema.org/a/6ksjrq7404k4877thbv5ega9r)

# How does it work ?

Instagram web is built on top of Django. Using [djangobot](https://github.com/hemantasapkota/djangobot), Insta logs you in as a web user. From this point onwards, it behaves exactly like a web browser. Note: I created djangobot as well.

**Does Instagram know about this ?** 

Pretty sure they do. If they wanted prevent something like this, they'd have re-captcha enabled on the login page, and not just the sign-up page. 

# Features

* Interactive console
* Manage multiple accounts
* Database log

# Roadmap

* True automation
* Console app
* Mobile apps ( IOS & Android )
* Chrome extension


