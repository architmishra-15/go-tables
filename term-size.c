#include <stdio.h>
#include <stdlib.h>

typedef struct {
    int width;
    int height;
} TerminalSize;


#if defined(WIN32) || defined(_WIN32) || defined(_WIN64) || defined(__WIN32__) || defined(__NT__)

#include <Windows.h>


TerminalSize* get_term_size() {
    TerminalSize* result = malloc(sizeof(TerminalSize));

    if (!result) return NULL;

    CONSOLE_SCREEN_BUFFER_INFO csbi;

    if (GetConsoleScreenBufferInfo(GetStdHandle(STD_OUTPUT_HANDLE), &csbi)){
        result->width = csbi.srWindow.Right - csbi.srWindow.Left + 1;
        result->height = csbi.srWindow.Bottom - csbi.srWindow.Top + 1;
    } else {
        free(result);
        return NULL;
    }

    return result;
}

#else

#include <unistd.h>
#include <sys/ioctl.h>


TerminalSize* get_term_size() {
    TerminalSize* result = malloc(sizeof(TerminalSize));
    if (!result) return NULL;
    
    struct winsize w;
    if (ioctl(STDOUT_FILENO, TIOCGWINSZ, &w) == 0) {
        result->width = w.ws_col;
        result->height = w.ws_row;
    } else {
        free(result);
        return NULL;
    }
    return result;
}

#endif

void free_terminal_size(TerminalSize* size) {
    if (size) free(size);
}
