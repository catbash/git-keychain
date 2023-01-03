#! /bin/bash

route=$HOME/.gitconfig

printf "
Select account to use.
  0: test (github.com)
  1: example (github.com)
>> "
read account
rm -drf $route

if [[ $account -eq 0 ]]
then
printf "
[user]
        name = test
        email = test@users.noreply.github.com
" >> $route
elif [[ $account -eq 1 ]]
then
printf "
[user]
        name = example
        email = example@users.noreply.github.com
" >> $route
else
echo "ERROR: Please select an account."
fi

echo "
Active git credentials:"
git config --list