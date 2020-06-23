# Declare known variables
IP_ADDRESS="192.168.178.98"
FILE_TO_SEND="Zentrale.salden"
USERNAME="ubuntu"

# Calls SFTP one-liner 
# Source: https://stackoverflow.com/questions/16721891/single-line-sftp-from-terminal
sftp $USERNAME@$IP_ADDRESS:$FILE_TO_SEND $FILE_TO_SEND