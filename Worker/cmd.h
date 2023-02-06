    MSG_BEGIN(REQ_TEST, 10 , stReqTesat){};
    BYTE type = 0; // xx
    DWORD pos = 0;   // xxx
    char playerName[MAX_USERNAME+1] = {0};
    MSG_END

    MSG_BEGIN(REQ_TEST2, 2, stReqTest2){};
    BYTE type = 0; // xx
    BYTE pos = 0;   // xxx
    MSG_END

    const BYTE SEND_IDENTIFY_RUNE_SC = 247;
    struct stSendIdentifyRune_SC : public stFreeMsg_1
    {
        stSendIdentifyRune_SC()
        {
            byParam = SEND_IDENTIFY_RUNE_SC;
            type = 0;
            level = 0;
        }
        BYTE type;
        DWORD level;
    };