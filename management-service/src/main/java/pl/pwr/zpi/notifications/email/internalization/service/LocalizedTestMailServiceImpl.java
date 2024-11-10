package pl.pwr.zpi.notifications.email.internalization.service;

import org.springframework.context.MessageSource;
import org.springframework.stereotype.Service;

@Service
public class LocalizedTestMailServiceImpl extends LocalizedMessageServiceImpl {

    public LocalizedTestMailServiceImpl(MessageSource testMailSource) {
        super(testMailSource);
    }
}
