#! /bin/tclsh
package require Tcl 8.6

set LOCAL_DIR [file dirname [file normalize [info script]]]

proc @comment {comment} {}

@comment {
	TODO;
	- check that the HOME variable is set;
	  otherwise throw error
	- other error handling
	- cleanup
	- how to run tcl on windows platforms?
	- test script on windows?
}

# initialize empty list to store accounts by name
set accounts [list]

proc {git account} {username email host {sshkey ""} {note ""}} {
	# initialize an account given a username, email, and host 
	# and add to global accounts list
	
	# set a prefix for the account array so that no variable
	# names get overwritten
	set {account name} "git__$username"

	# set up the account
	upvar ${account name} account
	set account(username) $username
	set account(email) $email
	set account(host) $host
	# ssh key defaults to username unless otherwise specified
	set account(sshkey) [expr {$sshkey eq "" ? $username : $sshkey}]
	set account(note) $note

	# append the account to the accounts list
	lappend ::accounts ${account name}
}

proc {git account username} {acc {upvar_lvl 1}} {
	# get an account username given an account array
	upvar $upvar_lvl $acc account
	return $account(username)
}

proc {print array} {a {prefix ""} {upvar_lvl 1}} {
	# print an array
	upvar $upvar_lvl $a arr
	foreach key [array names arr] {
		puts "$prefix$key: $arr($key)"
	}
}

proc {print accounts} {} {
	# todo: selective account print (i.e., don't print sshkey if
	# standard name, don't print note if not available, print
	# note at top of account info block if available
	for {set i 0} {$i < [llength $::accounts]} {incr i} {
		set account [lindex $::accounts $i]
		puts "$i: [{git account username} $account 2]"
		{print array} $account "\t" 2
	}
}

# import accounts from accounts.conf file
source [file nativename "$LOCAL_DIR/accounts.conf"]

array set LOG_TAGS {
	INFO {[ INFO ]}
	ERROR {[ ERROR ]}
}

proc {write config} {id} {
	# write ~/.ssh/config file
	set accountArr [lindex $::accounts $id]
	upvar 2 $accountArr account
	set config "Host $account(host)
	HostName $account(host)
	User git
	IdentityFile $::HOME/.ssh/$account(sshkey)
	IdentitiesOnly yes"
	set filename [file nativename "$::HOME/.ssh/config"]
	set {config file} [open $filename w]
	puts ${config file} $config
	close ${config file}
}

proc {write gitconfig} {id} {
	# write ~/.gitconfig file
	puts "$::LOG_TAGS(INFO) set to account $id"
	set accountArr [lindex $::accounts $id]
	upvar $accountArr account
	set gitconfig "\[user\]\n\tname = $account(username)\n\temail = $account(email)"
	set filename [file nativename "$::HOME/.gitconfig"]
	set {gitconfig file} [open $filename w]
	puts ${gitconfig file} $gitconfig
	close ${gitconfig file}
	{write config} $id
	puts [exec git config {--list}]
}


proc {validate account selection} {selection} {
	# determine if the account selection is
	# valid

	if {![string is integer $selection]} {
		# return false if the selection
		# is not an integer
		return false
	}
	if {[lindex $::accounts $selection] eq ""} {
		# return false if the selection
		# does not reference an account
		# index
		return false
	}
	if {[string length $selection] == 0} {
		# return false if the selection
		# is an empty string
		return false
	}
	return true
}

if {[lindex $argv 0] eq {list}} {
	{print accounts}
}

if {[lindex $argv 0] eq {set}} {
	set accID [lindex $argv 1]
	if {[{validate account selection} $accID]} {
		{write gitconfig} $accID
	} else {
		{print accounts}
		puts "select account:"
		gets stdin accID
		if {$accID eq ""} {
			puts "$LOG_TAGS(INFO) no selection made - exiting"
		} else {
			if {![{validate account selection} $accID]} {
				puts "$LOG_TAGS(ERROR) invalid selection - please try again"
			} else {
				puts "$LOG_TAGS(INFO) selected account: $accID"
				{print array} [lindex $::accounts $accID] "\t"
				{write gitconfig} $accID
			}
		}
	}
}

if {!$argc || [lindex $argv 0] eq {help}} {
	puts "$argv0
a cli written in tcl to manage git accounts across
multiple platforms

commands:
help 		display this help text
list		list all available accounts
set <index>	list all available accounts and then set account by index
		(you may optionally specify the account index to skip
		listing all accounts)
	"
}
