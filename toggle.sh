#! /bin/bash

route=/home/test/.gitconfig # set as route to your .gitconfig file

printf "
Select account to use.
  0: test
  1: example
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
