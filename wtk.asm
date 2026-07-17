format ELF64

	extrn   mfv_main

public main

SYS_WRITE       equ     1
SYS_EXIT        equ     60

macro write fd, buffer, count
	{
	mov     rax, SYS_WRITE
	mov     rdi, fd
	mov     rsi, buffer
	mov     rdx, count
	syscall
	}

macro exit code
	{
	mov     rax, SYS_EXIT
	mov     rdi, code
	syscall
	}

macro comp string, count
	{
	push    rsi
	push    rdi

	mov     rsi, [rsi + 8]
	lea     rdi, [string]
	mov     rcx, count

	repe    cmpsb
	}

macro prepreg offst
	{
	pop     rdi
	pop     rsi

	add     rsi, offst
	dec     edi
	sub     rsp, offst
	}

section '.text' executable

main:
	cmp     edi, 1
	jle     .err

	comp    str_mfv, mfv_len
	je      .run_mfv
	pop     rdi
	pop     rsi

	comp    str_help, help_len
	je      .help

	exit    0

.help:
	write   1, help_msg, help_msg_len
	exit    0

.err:
	write   1, error_msg, error_msg_len
	exit    1

.run_mfv:
	prepreg 8
	call    mfv_main

	add     rsp, 8
	exit    0


section '.data'
str_mfv         db      "mfv", 0
mfv_len =				$ - str_mfv

str_help        db      "--help", 0
help_len = 			$ - str_help

error_msg       db      "Not enough arguments were provided! See --help for usage", 10
error_msg_len = $ - error_msg

help_msg        db			"Usage: wtk <program_name> [args...]", 10
        				db      "Programs: ", 10
								db      " mfv - Missing Files Verifier", 10
help_msg_len = 	$ - help_msg

section '.note.GNU-stack'