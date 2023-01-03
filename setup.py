import json
import os
import argparse

SSH_CONFIG_FILE = f"/home/{os.environ.get('USER')}/.ssh/config" # active config file path
ENCODING = "utf-8"
BOLD_START = "\033[1m"
BOLD_END = "\033[0m"

TMP_EMAILS = [
    "test@users.noreply.github.com",
    "example@users.noreply.github.com"
]

def clear_file(file):
    """ Removes a file if it exists """
    if os.path.exists(file):
        os.remove(file)

def write_to_accounts_file(data):
    """ Given a list of dictionaries, overwrite to
    accounts.json file """
    clear_file("accounts.json")
    output = open("accounts.json", "w", encoding=ENCODING)
    with output as json_file:
        json_file.write(json.dumps(data))
    output.close()

def get_accounts_file_contents():
    """ Return contents of accounts.json as json data """
    with open('accounts.json') as json_file:
        return json.load(json_file)

def list_accounts():
    """ List accounts.json data """
    data = get_accounts_file_contents()
    print(f"\nNumber of accounts: {len(data)}\n")
    for entry in data:
        entry_info = (
            f"Name\t\t{entry['name']}\n"
            f"Email\t\t{entry['email']}\n"
            f"Platform\t{entry['platform']}"
        )
        if 'filename' in list(entry.keys()):
            entry_info += f"\nFilename\t{entry['filename']}"
        entry_info += "\n"
        print(entry_info)

def add_account(name, email, platform, filename):
    """ Add an account """
    if not name or not email or not platform:
        return 0
    data = get_accounts_file_contents()
    accounts = []
    for entry in data:
        if not entry['email'] in TMP_EMAILS:
            if (entry['email'] == email and entry['platform'] == platform) \
                or (entry['name'] == name and entry['platform'] == platform):
                print(f"\nERROR: Cannot add account. Duplicate name and/or email for given platform! Existing account: {entry}")
                return 0
            accounts.append(entry)

    new_account = {
        "name": name,
        "email": email,
        "platform": platform
    }

    if filename:
        new_account['filename'] = filename

    accounts.append(new_account)
    write_to_accounts_file(accounts)
    print("INFO: Account added!")

def delete_account(name, email, platform, filename):
    """ Delete an account """
    if (not name and not platform) or \
        (not email and not platform):
        return 0
    data = get_accounts_file_contents()
    filter_by = {
        'name': name,
        'email': email,
        'platform': platform,
        'filename': filename
    }
    found_account = False
    accounts = []
    for entry in data:
        modified_entry = {}
        tmp = {}
        for f,v in filter_by.items():
            if v:
                if f in list(entry.keys()):
                    modified_entry[f] = entry[f]
                tmp[f] = v
        if modified_entry == tmp:
            found_account = True
        else:
            accounts.append(entry)
    if not found_account:
        print("\nINFO: Account not found.")
        return 0
    write_to_accounts_file(accounts)
    print("\nINFO: Account deleted!")

def reset_accounts_file():
    """ Resets the accounts.json file if needed """
    data = [
      {
        "name": "test",
        "email": "test@users.noreply.github.com",
        "platform": "github.com"
      },
      {
        "name": "example",
        "email": "example@users.noreply.github.com",
        "platform": "github.com"
      },
    ]
    clear_file("accounts.json") # clear accounts file
    clear_file("config") # clear local config file
    write_to_accounts_file(data)

def generate_selection(data):
    selection = "\nSelect account to use.\n"
    for d in data:
        selection += f'  {data.index(d)}: {d["name"]} ({d["platform"]})\n'
    return selection

def generate_users(data):
    users = []
    for d in data:
        users.append(
            f'[user]\n        name = {d["name"]}\n        email = {d["email"]}'
        )
    return users

def generate_blocks(users):
    blocks = ""
    for user in users:
        if users.index(user) == 0:
            blocks += f'if [[ $account -eq {users.index(user)} ]]\nthen\nprintf "\n{user}\n" >> $route\n'
        else:
            blocks += f'elif [[ $account -eq {users.index(user)} ]]\nthen\nprintf "\n{user}\n" >> $route\n'
    blocks += 'else\necho "ERROR: Please select an account."\nfi'

    return blocks

def generate_ssh_config(data, file):
    """ Generates the config from provided data and writes to specified file """
    platforms = []
    clear_file(file)
    config = open(file, "w", encoding=ENCODING)
    for entry in data:
        host = entry['platform']
        id_file = f'/home/{os.environ.get("USER")}/.ssh/{entry["name"]}'
        if 'filename' in list(entry.keys()):
            id_file = f'/home/{os.environ.get("USER")}/.ssh/{entry["filename"]}'
        if host in platforms:
            host = entry['name']
        config.write(
            f"Host {host}\n"
            f"  HostName {entry['platform']}\n"
            f"  User git\n"
            f"  IdentityFile {id_file}\n"
            f"  IdentitiesOnly yes\n\n"
        )
        platforms.append(entry['platform'])
    config.close()

def generate_toggle_script():
    """ Generates the toggle script based on the accounts.json data """
    data = get_accounts_file_contents()

    route = f'/home/{os.environ.get("USER")}/.gitconfig'
    selection = generate_selection(data)
    users = generate_users(data)

    clear_file("toggle.sh")
    toggle = open("toggle.sh", "w", encoding=ENCODING)
    toggle.write(
        f'#! /bin/bash\n\nroute=$HOME/.gitconfig'
        f'\n\nprintf "{selection}>> "\nread account\nrm -drf $route\n\n'
        f'{generate_blocks(users)}\n\n'
        f'echo "\nActive git credentials:"\ngit config --list'
    )
    toggle.close()
    # generate live config file in /home/$USER/.ssh
    generate_ssh_config(data, SSH_CONFIG_FILE)
    # generate local copy of config file
    generate_ssh_config(data, "config")

def main():
    parser = argparse.ArgumentParser(description='Generates a script to toggle between multiple git accounts.')
    parser.add_argument("-r", "--reset", help="reset accounts.json template to example", action='store_true')
    parser.add_argument("-y", "--auto-confirm", help="skip prompt to confirm setup", action='store_true')
    parser.add_argument("-l", "--list", help="list configured accounts", action='store_true')
    parser.add_argument("-a", "--add", help="use this flag to add a new account; must specify name, email, and platform", action='store_true')
    parser.add_argument("-d", "--delete", help="use this flag to delete an account; must specify name and platform or email and platform", action='store_true')
    parser.add_argument("-n", "--name", help="specify name for adding/deleting an account", action='store', type=str)
    parser.add_argument("-e", "--email", help="specify email for adding/deleting an account", action='store', type=str)
    parser.add_argument("-p", "--platform", help="specify platform for adding/deleting an account", action='store', type=str)
    parser.add_argument("-f", "--filename", help="specify filename for adding/deleting an account", action='store', type=str)
    args = parser.parse_args()
    if args.reset:
        reset_accounts_file()
        print("INFO: Reset successful!")
    elif args.delete:
        if (not args.email and not args.platform) \
            or (not args.name and not args.platform):
            print(f"ERROR: Please supply a name and platform or email and platform to identify the account you want to delete. Example commands:\n  python setup.py -d -name=test -p=github.com\n  python setup.py -d -e=test@test.com -p=github.com")
        delete_account(args.name, args.email, args.platform, args.filename)
    elif args.add:
        if not args.email or not args.name or not args.platform:
            print(f"ERROR: Please supply a name, email, and platform for the new account! Example command:\n python setup.py -a -e=example@example.com -n=example -p=github.com")
        add_account(args.name, args.email, args.platform, args.filename)
    elif args.list:
        list_accounts()
    else:
        proceed = False
        if args.auto_confirm:
            proceed = True
        else:
            print("This script will overwrite your /home/$USER/.ssh/config file. Are you sure you want to continue?")
            confirm = input("Type YES to confirm. >> ")
            if confirm == "YES":
                proceed = True
        if proceed:
            generate_toggle_script()
            print(f"\nINFO: Script generated! To complete setup, add an alias for the following:")
            print(f"\t{BOLD_START}{os.environ.get('PWD')}/toggle.sh{BOLD_END}")
        else:
            print("\nINFO: Setup aborted!")

if __name__=="__main__":
    main()
