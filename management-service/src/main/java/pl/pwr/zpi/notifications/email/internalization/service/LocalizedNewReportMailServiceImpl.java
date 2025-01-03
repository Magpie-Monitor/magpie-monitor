package pl.pwr.zpi.notifications.email.internalization.service;

import org.springframework.context.MessageSource;
import org.springframework.stereotype.Service;

@Service
public class LocalizedNewReportMailServiceImpl extends LocalizedMessageServiceImpl {

	public LocalizedNewReportMailServiceImpl(MessageSource newReportMailSource) {
		super(newReportMailSource);
	}
}
