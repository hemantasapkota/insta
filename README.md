# Insta - Power tool for Instagram.

Insta is a command line power tool for Instagram. With it, you can like, follow, comment and query Instagram.

# Motivation

There're plenty of Instagram automation tools floating around the web and app stores. Most of them, however, ask for your username and password which poses a big safety risk.

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

# How does it work ?

Instagram web is built on top of Django. Using [djangobot](https://github.com/hemantasapkota/djangobot), Insta logs you in as a web user. From this point onwards, it behaves exactly like a web browser.

# Features

* Interactive console
* Manage multiple accounts
* Database log
* Built-in DSL for scripting flows

# Flows

Insta comes with a built-in DSL that can be used to script flows. Example flow listed below:

```
$(get_data path=explore/tags/coffee => result)
$(loop range=0,5 => i)
  $(filter var=result query=entry_data.TagPage[0].tag.media.nodes[$(`i`)].display_src => url)
  $(download url="$(`url`)")
$(pool)

```

The script queries for the tag *coffee*, loops through the result 5 times, filters the image url from the result and downloads it.

# Roadmap

* True automation
* CLI app
* Mobile apps ( IOS & Android )
* Chrome extension


