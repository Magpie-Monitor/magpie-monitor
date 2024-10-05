package pl.pwr.zpi.notifications.email.html.service;


import org.commonmark.node.Node;
import org.commonmark.parser.Parser;
import org.commonmark.renderer.html.HtmlRenderer;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.notifications.email.html.styling.HeadingAttributeProvider;
import pl.pwr.zpi.notifications.email.html.styling.TextAttributeProvider;

@Service
public class MarkdownServiceImpl implements MarkdownService {

    private final Parser parser;
    private final HtmlRenderer renderer;

    public MarkdownServiceImpl() {
        this.parser = Parser.builder().build();

        this.renderer = HtmlRenderer.builder()
                .attributeProviderFactory(attributeProviderContext -> new HeadingAttributeProvider())
                .attributeProviderFactory(attributeProviderContext -> new TextAttributeProvider())
                .build();
    }

    @Override
    public String toHtmlWithMarkdowns(String markdownText) {
        Node document = parser.parse(markdownText);
        return renderer.render(document);
    }
}
