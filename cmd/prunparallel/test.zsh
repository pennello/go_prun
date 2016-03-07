echo index "$@"
sleep $(($RANDOM % 2 + 1))

# Random exit code.
r=
if [ $(($RANDOM % 4)) = 0 ];then
  r=$(($RANDOM % 32 + 1))
else
  r=0
fi
echo exiting with $r
exit $r
