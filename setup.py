import json
import os
import argparse

ENCODING = "utf-8"

def clear_file(file):
    """ Removes a file if it exists """
    if os.path.exists(file):
        os.remove(file)

def reset_accounts_file():
    """ Resets the accounts.json file if needed """
    data = [
      {
        "name": "test",
        "email": "test@users.noreply.github.com"
      },
      {
        "name": "example",
        "email": "example@users.noreply.github.com"
      },
    ]
    clear_file("accounts.json")
    output = open("accounts.json", "w", encoding=ENCODING)
    with output as json_file:
        json_file.write(json.dumps(data))
    output.close()

def generate_selection(data):
    selection = "\nSelect account to use.\n"
    for d in data:
        selection += f'  {data.index(d)}: {d["name"]}\n'
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


def generate_toggle_script():
    """ Generates the toggle script based on the accounts.json data """
    with open('accounts.json') as json_file:
        data = json.load(json_file)

    route = f'/home/{os.environ.get("USER")}/.gitconfig'
    selection = generate_selection(data)
    users = generate_users(data)

    clear_file("toggle.sh")
    toggle = open("toggle.sh", "w", encoding=ENCODING)
    toggle.write(
        f'#! /bin/bash\n\nroute={route} # set as route to your .gitconfig file'
        f'\n\nprintf "{selection}>> "\nread account\nrm -drf $route\n\n'
        f'{generate_blocks(users)}\n\n'
        f'echo "\nActive git credentials:"\ngit config --list'
    )
    toggle.close()

def main():
    parser = argparse.ArgumentParser(description='Generates a script to toggle between multiple git accounts.')
    parser.add_argument("-r", "--reset", help="reset accounts.json template", action='store_true')
    args = parser.parse_args()
    if args.reset:
        reset_accounts_file()
    else:
        generate_toggle_script()

if __name__=="__main__":
    main()
