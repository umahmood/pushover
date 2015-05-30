# Pushover

Pushover is an easy to use command line tool for sending messages via the great
Pushover.net service. Pushover.net is a platform for sending and receiving push 
notifications, see here for more details https://pushover.net/ .

**Note**: This project is not affiliated with Pushover.net or Superblock LLC in 
any way.

# Installation

You will need to set up and configure the ![Go](https://golang.org/doc/install) 
tool chain.

> go get github.com/umahmood/pushover
> cd $GOPATH/src/github.com/umahmood/pushover
> go install

The binary should now be located in: $GOTPATH/bin.

# Setup

##### 1. Find your token and user keys:

In order to send messages, the pushover command line tool needs to read your 
token and user keys. You can find these key by logging into your Pushover 
account. To find your user key, it will be on the home dashboard under 
'Your User Key':

![ukey](https://i.imgur.com/eRcPDgx.png)

To find your token key, click on an application in the 'Your Applications' list:

![tkey](https://i.imgur.com/TphZMpo.png)

##### 2. Create a dot file:

We now need to create a dot file in your home directory, containing your token
and user keys. This will stop you having to enter them every time the tool is 
invoked.

On Linux/OSX open a terminal and enter:

> echo -e "token=YYY\nuser=XXX" >> ~/.pushover

Replace YYY and XXX with your token and user keys respectively.

On Windows:

1. Go to your home directory pointed to by the %UserProfile% environment variable, 
for example on Windows 7 this is C:/Users/your_user_name.

2. Create the dot file .pushover ([see here](https://gist.github.com/ozh/4131243) on how to create a dot file).

3. Open the file with Notepad, enter and then save the content:
> token=YYY<br/>user=XXX

Replace YYY and XXX with your token and user keys respectively.

# Usage

> pushover -msg "hello world"

> pushover -title "Hi" -msg "Hello World" -device "nexus5" -priority 2 -timestamp 1433024710 -url "http://www.google.com" -url-title "Google"

> pushver -h <br>
Usage of pushover:<br>
  -device="": Send message directly to this device, rather than all devices.<br>
  -msg="": (Required) - Your message.<br>
  -priority="": Message priority. -2 = lowest , -1 = low, 0 = normal, 1 = high, 2 = emergency.<br>
  -sound="": Sound to play when user receives notification, overrides the user's default sound choice.<br>
  -timestamp="": A Unix timestamp of your message's date and time to display to the user.<br>
  -title="": Your message's title, otherwise your app's name is used.<br>
  -url="": A supplementary URL to show with your message.<br>
  -url-title="": A title for your supplementary URL, otherwise just the URL is shown.<br>
  -v=false: Display message response details.<br>

# License

See the [LICENSE](LICENSE.md) file for license rights and limitations (MIT).
