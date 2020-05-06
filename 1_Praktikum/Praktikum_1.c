#include <unistd.h>
#include <sys/wait.h>
#include <stdlib.h>
#include <sys/stat.h>		// open()
#include <signal.h>			// signal()
#include <fcntl.h>			// open()
#include <stdio.h>			// printf(), ...
#include <time.h>				// time(), ...
#include <string.h>			// strtok()
#include <stdbool.h>
#define MAXLINE 100
#define MOD "exit with CTR C"
#define clear() printf("\033[H\033[J")

void init_shell()
{
    clear();
    printf("\n******************"
           "************************");
    printf("\n\t******MSH*******");
    printf("\n*******************"
           "***********************");
    clear();
}

void printDirectory()
{
    char cwd[1024];
    getcwd(cwd, sizeof(cwd));
    printf("\n%s>", cwd);
}

void sig_handler(int sig)
{
    printf("received SIGINT\n");
    exit(1);
}

static void sigchld_hdl ( int sig ) {
    while ( waitpid ( -1 , NULL , WNOHANG ) > 0);
}

int read_command(char *command, char *parameters[]) { // prompt for user input and read a command line
    printDirectory();
    int noParam = 0;
    size_t bufsize = MAXLINE;
    ssize_t characters;

    if( command == NULL)
    {
        perror("Unable to allocate buffer");
        exit(1);
    }
    characters = getline(&command,&bufsize,stdin);
    if ((command)[characters - 1] == '\n')
    {
        (command)[characters - 1] = '\0';
        --characters;
    }

    //getting the number of arguments
    int count =0;
    int i;

    for(i=0;i<strlen(command);i++){
        if(command[i] == ' ')
            count++;
    }
    count++;
    //split command by space
    char* parts = strtok(command, " ");
    //start parameters with the command
    parameters[0] = command;

    //fill parameters
    for (int i = 0; i < count; i++)
    {
        parameters[i] = parts;
        parts = strtok(NULL, " ");
    }
    //insert NULL after all other parameters
    parameters[count] = NULL;
    return(count);
} // read_command

//return 0 if last char = & else return 1
int background(char* command) {
    return (command && *command && command[strlen(command) - 1] == '&') ? 0 : 1;
}

int main(int argc, char *argv[]) {
    int childPid;
    int status;
    char command[MAXLINE];
    char *parameters[MAXLINE];
    int noParams;
    int back;
    signal(SIGINT, sig_handler);
    if (signal(SIGCHLD, sigchld_hdl)) {
        perror ("signal");
        return 1;
    }

    init_shell();

    while (1) {
        noParams = read_command(command, parameters); // read user input
        back = background(parameters[noParams-1]);
        //if background execution remove & from parameters
        if (back == 0)
        {
            for (int i = 0; i < noParams; i ++)
            {
               if (strcmp(parameters[i],"&") == 0)
               {
                   parameters[i] = NULL;
               }
            }
        }
        if (noParams == 0)
        {
            fprintf(stderr, "no command ?!\n");
            exit(1);
        }
        if ((childPid = fork()) == -1)
        { // create process
            fprintf(stderr, "can't fork!\n");
            exit(2);
        }
        //check if the command is cd
        else if (strcmp("cd", command) == 0)
        {
            //Wenn cd ohne parameter ausgeführt wird wechsle in das Home verzeichnis
            if (parameters[1] == NULL)
            {
                char* username;
                username = getlogin();
                char buf[256];
                snprintf(buf, sizeof buf, "/home/%s", username);
                chdir(buf);
            }
            chdir(parameters[1]);
        }
        //check if the command is quit or exit
        else if (strcmp("quit", command) == 0 || strcmp("exit", command) == 0)
        {
            exit(1);
        }
        else if (childPid == 0) { // child process
            execvp(command, parameters);
                perror("execution failed");
            exit(3);
        }
        //search for a & to decidde if we should wait or nah
        else
        {
            if (back == 1)
            {
                waitpid(childPid, &status, WUNTRACED | WCONTINUED);
            }
        }
    }
    exit(0);
}