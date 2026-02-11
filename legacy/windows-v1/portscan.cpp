#include <iostream>
#include <WinSock2.h>
#include <WS2tcpip.h>
#include <thread> 
#include <vector>
#include <fstream>
#include <string>
#include <wininet.h>

#pragma comment(lib, "wininet.lib")
#pragma comment(lib, "ws2_32.lib")

using namespace std;

void HandleTCPConnection(std::string destination_ip, unsigned short start_port, unsigned short end_port, int max_attempts, std::ostream& output)
{
    WSADATA wsa_data;
    if (WSAStartup(MAKEWORD(2, 2), &wsa_data) != 0)
    {
        output << "Error initializing Winsock" << endl;
        return;
    }

    for (unsigned short port = start_port; port <= end_port; port++)
    {
        int attempts = 0;
        SOCKET socket_fd;
        cout << "Scanning port " << port << "..." << endl;

        while (attempts < max_attempts)
        {
            socket_fd = socket(AF_INET, SOCK_STREAM, IPPROTO_TCP);
            if (socket_fd == INVALID_SOCKET)
            {
                output << "Error creating socket: " << WSAGetLastError() << endl;
                WSACleanup();
                return;
            }

            sockaddr_in dest_addr;
            memset(&dest_addr, 0, sizeof(dest_addr));
            dest_addr.sin_family = AF_INET;
            dest_addr.sin_port = htons(port);

            wchar_t dest_ip_wide[16];
            MultiByteToWideChar(CP_UTF8, 0, destination_ip.c_str(), -1, dest_ip_wide, 16);
            InetPton(AF_INET, dest_ip_wide, &(dest_addr.sin_addr));

            if (connect(socket_fd, (struct sockaddr*)&dest_addr, sizeof(dest_addr)) == SOCKET_ERROR)
            {
                output << "Error connecting to " << destination_ip << ":" << port << " - Error code: " << WSAGetLastError() << endl;
                closesocket(socket_fd);
                attempts++;
            }
            else
            {
                output << destination_ip << ":" << port << " is open." << endl;
                closesocket(socket_fd);
                WSACleanup();
                return;
            }
        }

        output << "Failed to connect to " << destination_ip << ":" << port << " after " << max_attempts << " attempts." << endl;
    }

    WSACleanup();
    return;
}

void HandleUDPConnection(std::string destination_ip, unsigned short start_port, unsigned short end_port, int max_attempts, std::ostream& output)
{
    WSADATA wsa_data;
    if (WSAStartup(MAKEWORD(2, 2), &wsa_data) != 0)
    {
        output << "Error initializing Winsock" << endl;
        return;
    }

    for (unsigned short port = start_port; port <= end_port; port++)
    {
        int attempts = 0;
        SOCKET socket_fd;
        cout << "Scanning port " << port << " (UDP)..." << endl;

        while (attempts < max_attempts)
        {
            socket_fd = socket(AF_INET, SOCK_DGRAM, IPPROTO_UDP);
            if (socket_fd == INVALID_SOCKET)
            {
                output << "Error creating socket: " << WSAGetLastError() << endl;
                WSACleanup();
                return;
            }

            sockaddr_in dest_addr;
            memset(&dest_addr, 0, sizeof(dest_addr));
            dest_addr.sin_family = AF_INET;
            dest_addr.sin_port = htons(port);

            wchar_t dest_ip_wide[16];
            MultiByteToWideChar(CP_UTF8, 0, destination_ip.c_str(), -1, dest_ip_wide, 16);
            InetPton(AF_INET, dest_ip_wide, &(dest_addr.sin_addr));

            int send_result = sendto(socket_fd, "", 1, 0, (struct sockaddr*)&dest_addr, sizeof(dest_addr));

            if (send_result == SOCKET_ERROR) {
                output << "Error sending to " << destination_ip << ":" << port << " (UDP) - Error code: " << WSAGetLastError() << endl;
            }
            else {
                output << destination_ip << ":" << port << " is open (UDP)." << endl;
            }

            closesocket(socket_fd);
            WSACleanup();
            return;
        }

        output << "Failed to connect to " << destination_ip << ":" << port << " after " << max_attempts << " attempts." << endl;
    }

    WSACleanup();
    return;
}

void HandleFTPConnection(std::string destination_ip, int max_attempts, std::ostream& output)
{
    HINTERNET hInternet = InternetOpen(L"FTP Checker", INTERNET_OPEN_TYPE_DIRECT, NULL, NULL, 0);
    if (!hInternet)
    {
        output << "Error initializing WinINet" << endl;
        return;
    }

    wstring wideString = wstring(destination_ip.begin(), destination_ip.end());
    LPCWSTR wideChar = wideString.c_str();


    HINTERNET hFtpSession = InternetConnect(hInternet, wideChar, INTERNET_DEFAULT_FTP_PORT, L"", L"", INTERNET_SERVICE_FTP, 0, 0);
    if (!hFtpSession)
    {
        output << "Error connecting to FTP server - Error code: " << GetLastError() << endl;
        InternetCloseHandle(hInternet);
        return;
    }

    output << destination_ip << " FTP port is open." << endl;

    InternetCloseHandle(hFtpSession);
    InternetCloseHandle(hInternet);
}

int main()
{
    int max_attempts = 1;
    unsigned short start_port, end_port;
    int scan_type;

    cout << "Enter scan type (1 for TCP, 2 for UDP, 3 for FTP): ";
    cin >> scan_type;

    switch (scan_type)
    {
    case 1:
        cout << "Enter the starting port: ";
        cin >> start_port;
        cout << "Enter the ending port: ";
        cin >> end_port;

        while (true)
        {
            string destination_ip;
            cout << "Enter the destination IP address (or 'exit' to quit): ";
            cin >> destination_ip;

            if (destination_ip == "exit")
                break;

            thread t1(HandleTCPConnection, destination_ip, start_port, end_port, max_attempts, ref(cout));
            t1.join();
        }
        break;

    case 2:
        while (true)
        {
            cout << "Enter the starting port: ";
            cin >> start_port;
            cout << "Enter the ending port: ";
            cin >> end_port;

            string destination_ip;
            cout << "Enter the destination IP address (or 'exit' to quit): ";
            cin >> destination_ip;

            if (destination_ip == "exit")
                break;

            thread t1(HandleUDPConnection, destination_ip, start_port, end_port, max_attempts, ref(cout));
            t1.join();
        }
        break;

    case 3:
        while (true)
        {
            cout << "Enter the starting port: ";
            cin >> start_port;
            cout << "Enter the ending port: ";
            cin >> end_port;

            string destination_ip;
            cout << "Enter the destination IP address (or 'exit' to quit): ";
            cin >> destination_ip;

            if (destination_ip == "exit")
                break;

            thread t1(HandleFTPConnection, destination_ip, max_attempts, ref(cout));
            t1.join();
        }
        break;

    default:
        cout << "Invalid scan type. Please enter 1 for TCP, 2 for UDP, 3 for FTP." << endl;
    }

    return 0;
}
