# Get variables from the sftp.conf file
# Source: https://askubuntu.com/questions/743493/best-way-to-read-a-config-file-in-bash
. ./sftp.conf

# Calls SFTP one-liner 
# Source: https://stackoverflow.com/questions/16721891/single-line-sftp-from-terminal
sftp $USERNAME@$IP_ADDRESS:$FILE_TO_SEND "."