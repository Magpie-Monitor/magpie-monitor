for i in {1..1000000};
do
  echo "testwetwerwedwxsxwsxwsxwxw$i" >> test.log
#	echo $(tr -dc A-Za-z0-9 </dev/urandom | head -c 13) >> test.log
#	echo $(tr -dc A-Za-z0-9 </dev/urandom | head -c 13) >> test2.log
	sleep 1
done

