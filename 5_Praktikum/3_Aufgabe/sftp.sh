# Get variables from the sftp.conf file
# Source: https://askubuntu.com/questions/743493/best-way-to-read-a-config-file-in-bash
. ./sftp.conf

# Calls SFTP one-liner 
# Source: https://stackoverflow.com/questions/16721891/single-line-sftp-from-terminal
sftp $USERNAME@$IP_ADDRESS:$DIRECTORY"*" "."

num_of_lines=0;
betrag_sum=0;
# Here we read every line from the "Zentrale.salden" file and save it into $line
# Source: https://linuxhint.com/read_file_line_by_line_bash/
while read line; do
    # Here we convert the line into an array of strings divided by spaces
    # Source: https://stackoverflow.com/a/13402368
    line_array=($line)
    betrag="${line_array[1]}"
    if [[ "${betrag:0:1}" == "+" ]]; then
        # Here we first get the number of "betrag" with ${betrag:1}, then we
        #  convert it to decimal (Because the many "0" are giving the bash many problems)
        #  and then we accumulate it in $betrag_sum
        # Source if/else: https://linuxacademy.com/blog/linux/conditions-in-bash-scripting-if-statements/
        # Source adition: https://unix.stackexchange.com/questions/55069/how-to-add-arithmetic-variables-in-a-script
        # Source decimal: https://stackoverflow.com/questions/24777597/value-too-great-for-base-error-token-is-08
        num_of_lines=$(($num_of_lines+1))
        betrag_sum=$((betrag_sum + ((10#${betrag:1}))))
    elif [[ "${betrag:0:1}" == "-" ]]; then
        num_of_lines=$(($num_of_lines+1))
        betrag_sum=$((betrag_sum - ((10#${betrag:1}))))
    fi;
done < "Zentrale.salden";  

# Here we compare the calculated values with the values of the "Zentrale.pruefsumme" file
# For that we get each paragraph of that file and compare it to our values
# Source: https://www.geeksforgeeks.org/write-bash-script-print-particular-line-file/
if [ "$(sed -n 1p Zentrale.pruefsumme)" = "$num_of_lines" ]; then
    if [ "$(sed -n 2p Zentrale.pruefsumme)" = "$betrag_sum" ]; then
        echo "pass"
    else
        echo "Der Betrag stimmt nicht mit die Pruefsumme"
    fi
else
    echo "Die Nummer der Linien stimmt nicht mit die Pruefsumme"
fi

