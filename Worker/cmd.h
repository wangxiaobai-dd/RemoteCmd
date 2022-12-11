    MSG_BEGIN(REQ_TEST, 10 , stReqTest){};
    BYTE type = 0; // xx
    DWORD pos = 0;   // xxx
    char playerName[MAX_USERNAME+1] = {0};
    MSG_END

    MSG_BEGIN(REQ_TEST2, 2, stReqTest2){};
    BYTE type = 0; // xx
    BYTE pos = 0;   // xxx
    MSG_END
