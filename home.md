---
title: Magpie Monitor
description: 
published: true
date: 2024-12-15T12:06:17.460Z
tags: 
editor: markdown
dateCreated: 2024-12-02T23:31:18.691Z
---

# System raportowania zdarzeń systemowych w języku naturalnym z wykorzystaniem dużych modeli językowych

## Dokumentacja

**Spis treści**

[1\. Wykaz symboli, oznaczeń i akronimów](#wykaz-symboli-oznaczeń-i-akronimów)

[2\. Słownik pojęć](#słownik-pojęć)

[3\. Cel i zakres przedsięwzięcia](#cel-i-zakres-przedsięwzięcia)

[3.1 Cel](#cel)

[3.2 Zakres](#zakres)

[4\. Stan wiedzy w obszarze przedsięwzięcia](#stan-wiedzy-w-obszarze-przedsięwzięcia)

[5\. Założenia wstępne](#założenia-wstępne)

[5.1. Przeznaczenie](#przeznaczenie)

[5.2. Nazwa i logo projektu](#nazwa-i-logo-projektu)

[5.3. Analiza funkcjonalna](#analiza-funkcjonalna)

[5.4. Przyjęte ograniczenia](#przyjęte-ograniczenia)

[6. Specyfikacja i analiza wymagań na produkt programowy](#specyfikacja-i-analiza-wymagań-na-produkt-programowy)

[6.1. Użytkownicy systemu](#użytkownicy-systemu)

[6.2. Wymagania niefunkcjonalne](#wymagania-niefunkcjonalne)

[6.3. Wymagania funkcjonalne](#wymagania-funkcjonalne)

[6.4. Historyjki użytkownika](#historyjki-użytkownika)

[6.5. Mapowanie wymagań funkcjonalnych na historyjki użytkownika](#mapowanie-wymagań-funkcjonalnych-na-historyjki-użytkownika)

[7. Projekt produktu programowego](#projekt-produktu-programowego)

[7.1. Wykorzystane technologie i narzędzia](#wykorzystane-technologie-i-narzędzia)

[7.2. Architektura systemu](#architektura-systemu)

[7.2.1. Model C4 \- poziom 1](#model-c4---poziom-1)

[7.2.2. Model C4 \- poziom 2](#model-c4---poziom-2)

[7.2.3. Diagram wdrożenia](#diagram-wdrożenia)

[7.2.4. Model C3 \- poziom 3 \- Agent](#model-c3---poziom-3---agent)

[7.2.5. Model C3 \- poziom 3 \- Logs Ingestion Service](#model-c3---poziom-3---logs-ingestion-service)

[7.2.6. Model C3 \- poziom 3 \- Report Service](#model-c3---poziom-3---report-service)

[7.2.7. Model C3 \- poziom 3 \- Metadata Service](#model-c3---poziom-3---metadata-service)

[7.2.8. Model C3 \- poziom 3 \- Management API](#model-c3---poziom-3---management-api)

[7.2.9. Model C3 \- poziom 3 \- Client](#model-c3---poziom-3---client)

[7.3. Bazy danych mikroserwisów](#bazy-danych-mikroserwisów)

[7.3.1. Baza logów](#baza-logów)

[7.3.2. Baza raportów](#baza-danych-raportów)

[7.3.3. Bazy danych metadata service](#bazy-danych-metadata-service)

[7.3.4. Bazy danych management service](#bazy-danych-management-service)

[7.4. Interfejsy programistyczne](#interfejsy-programistyczne)

[7.5. Projekt interfejsu](#projekt-interfejsu)

[7.5.1. Widok logowania](#widok-logowania)

[7.5.2. Widok główny](#widok-główny)

[7.5.3. Widok incydentu](#widok-incydentu)

[7.5.4. Widok raportów](#widok-raportów)

[7.5.5. Widok klastrów](#widok-klastrów)

[7.5.6. Widok konfiguracji raportów](#widok-konfiguracji-raportów)

[7.5.7. Widok konfiguracji powiadomień](#widok-konfiguracji-powiadomień)

[7.6 Diagramy procesów](#diagramy-procesów)

[8. Implementacja](#implementacja)

[8.1. Środowisko pracy](#środowisko-pracy)

[8.2. Struktura plików projektu](#struktura-plikow-projektu)

[8.3. Struktura plików w aplikacji “Agenta”](#struktura-plików-w-aplikacji-“agenta”)

[8.4. Struktura plików w serwisie “Ingestion Service”](#struktura-plików-w-serwisie-“ingestion-service”)

[8.5. Struktura plików w serwisie “Report Service”](#struktura-plików-w-serwisie-“report-service”)

[8.6. Struktura plików w serwisie “Metadata Service”](#struktura-plików-w-serwisie-“metadata-service”)

[8.7. Struktura plików w aplikacji “Management Service”](#struktura-plików-w-aplikacji-“management-service”)

[8.8. Struktura plików w aplikacji klienckiej](#struktura-plików-w-aplikacji-klienckiej)

[8.9. Uwierzytelnienie użytkownika](#uwierzytelnienie-użytkownika)

[8.10. Planowanie raportów (scheduling raportów, management service)](<#planowanie-raportów-(scheduling-raportów,-management-service)>)

[8.11. Zbieranie logów (agent)](<#zbieranie-logów-(agent)>)

[8.12. Zapisywanie logów (ingestion service)](<#zapisywanie-logów-(ingestion-service)>)

[8.13. Generowanie raportów (reports service)](<#generowanie-raportów-(reports-service)>)

[8.14. Ustawianie kanałów komunikacji (management service)](<#ustawianie-kanałów-komunikacji-(management-service)>)

[8.15. Odczytywanie stanu klastra (metadata service)](<#odczytywanie-stanu-klastra-(metadata-service)>)

[8.16. Zabezpieczenia aplikacji (management service)](<#zabezpieczenia-aplikacji-(management-service)>)

[9\. Testy produktu programowego/Wyniki i analiza badań](#testy-produktu-programowego/wyniki-i-analiza-badań)https://wikijs.rolo-labs.xyz/e/en/home#testy-agenta

[9.1. Testy Reports Service](#testy-reports-service)

[9.2. Testy Logs Ingestion Service](#testy-logs-ingestion-service)

[9.3. Testy Agenta](#testy-agenta)

[9.4. Testy Metadata Service](#testy-metadata-service)

[9.5. Testy Management Service](#testy-management-service)

[9.6. Testy funkcjonalne](#testy-funkcjonalne)

[10. Podsumowanie](#podsumowanie)

[10.1. Przebieg projektu](#przebieg-projektu)

[10.2. Wnioski](#wnioski)

[11\. Dokumentacja użytkownika](#dokumentacja-użytkownika)

[11.1\. Wprowadzenie](#wprowadzenie)

[11.1.1 Użytkowanie produktu programowego](#użytkowanie-produktu-programowego)

[11.1.2. Instalacja aplikacji](#instalacja-aplikacji)

[11.2.2 Najczęściej wykonywane operacje](#najczęściej-wykonywane-operacje)

[11.2.2.1. Logowanie do aplikacji](#logowanie-do-aplikacji)

[11.2.2.2 Otworzenie ostatniego raportu](#otworzenie-ostatniego-raportu)

[11.2.2.3. Planowanie generowania raportów](#planowanie-generowania-raportów)

[11.2.2.4. Generacja raportu na życzenie](#generacja-raportu-na-życzenie)

[11.2.2.5. Konfiguracja kanałów powiadomień](#konfiguracja-kanałów-powiadomień)

# 1. Wykaz symboli, oznaczeń i akronimów {#wykaz-symboli,-oznaczeń-i-akronimów}

| Akronim        | Znaczenie                                                        |
| :------------- | :--------------------------------------------------------------- |
| Magpie Monitor | Nazwa zrealizowanego projektu                                    |
| OpenAI         | Dostawca modeli językowych wykorzystanych do realizacji projektu |

# 2. Słownik pojęć {#słownik-pojęć}

| Pojęcie    | Synonimy | Znaczenie                                                                                                                                                                                                  |
| :--------- | :------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Aplikacja  |          | Aplikacja w kontekście klastra Kubernetesa. Deployment albo Statefull Set.                                                                                                                                 |
| Host       |          | Urządzenie z systemem Linux                                                                                                                                                                                |
| Obserwacja | Insight  | Wykryty problem z logów aplikacji/hostów. Obserwacja zawiera zarówno źródła obserwacji (na podstawie czego obserwacja była wysunięta) oraz podsumowanie i rekomendacje rozwiązania danego problemu         |
| Incydent   | Incident | Obserwacja uzupełniona o dokładniejsze metadane z jej źródła. Zawiera między innymi czas wystąpienia incydentu, metadane hosta/aplikacji, z którego pochodzi oraz konfiguracje z jaką została wygenerowana |
| Pilność    | Urgency  | Używana w kontekście raportu i incydentu, określa jak bardzo dany problem (określany przez incydent lub raport) wymaga natychmiastowej naprawy                                                             |

# 3. Cel i zakres przedsięwzięcia {#cel-i-zakres-przedsięwzięcia}

## 3.1 Cel {#cel}

Celem projektu jest opracowanie kompleksowej, wysoko skalowalnej aplikacji webowej, wspieranej przez program instalowany na wybranym klastrze komputerowym. System będzie umożliwiał monitorowanie i analizę logów systemowych oraz generowanie szczegółowych raportów na ich podstawie. Opracowywane rozwiązanie ma dostarczyć administratorom zaawansowane narzędzie, które usprawni szybkie i efektywne zarządzanie dużymi wolumenami danych generowanych przez różnorodne aplikacje.

Projektowany system ma wspierać użytkowników w kluczowych obszarach, takich jak:

- identyfikacja potencjalnych problemów,
- rozwiązywanie zidentyfikowanych zagrożeń,
- podejmowanie decyzji opartych na przejrzystych i klarownych raportach.

W założeniu efektem projektu powinno być znaczące podniesienie jakości monitorowania aplikacji oraz uproszczenie procesów analizy i zarządzania danymi logów w sposób intuicyjny, skuteczny i dostosowany do potrzeb użytkowników.

Z uwagi na nowatorskie podejście do monitorowania aplikacji, dodatkowym celem projektu będzie ocena skuteczności przyjętego rozwiązania. Pozwoli to określić jego wartość biznesową w porównaniu do bardziej standardowych metod, a także wskazać potencjalne kierunki dalszego rozwoju systemu.

## 3.2 Zakres {#zakres}

Zakres projektu prezentuje się następująco:

- Badanie potrzeb rynku oraz dostępnych rozwiązań analizujących logi systemowe,
- Sformułowanie wymagań koniecznych do osiągnięcia zaplanowanego celu
- Zaprojektowanie skalowalnej architektury systemu, opartej na niezależnych od siebie mikroserwisach.
- Zaprojektowanie intuicyjnego oraz estetycznego wyglądu aplikacji
- Implementacja oraz wdrożenie systemu na chmurowej platformie Linode.
- Testowanie systemu, zarówno z użyciem testów automatycznych oraz manualnych.
- Analiza osiągniętych rezultatów względem tradycyjnych rozwiązań monitorowania logów.

# 4. Stan wiedzy w obszarze przedsięwzięcia {#stan-wiedzy-w-obszarze-przedsięwzięcia}

Aplikacja Magpie Monitor dotyczy branży obserwowania systemów (ang. observability), będącą poddziedziną szerszej dziedziny technik analizy oraz przetwarzania danych. Obszar observability skupia się w głównej mierze na analizie danych emitowanych przez aplikacje, a jego filarami są metryki, logi oraz smugi (ang. traces).

Domeną Magpie Monitor jest semantyczna analiza logów emitowanych przez klastry komputerowe oraz przedstawianie rekomendacji w przypadku wykrycia nieprawidłowości. Obszar ten nie cieszy się wysoką popularnością w branży, co potwierdza fakt, jako iż najpopularniejsze narzędzia takie jak Dynatrace [17], Datadog [18] czy Logz.io [19], nie dokonują takowej analizy, czego przyczyny można doszukiwać się w niedawnej ekspansji technologii dużych modeli językowych.

Znakomita większość narzędzi zorientowanych na logi emitowane z aplikacji, takich jak Azure Monitor [20] czy Amazon CloudWatch [21], jest skupiona na komfortowym wyświetlaniu oraz wyszukiwaniu logów, a także analizie wzorców, które w logach występują. Podejście to pozwala ekspertom domenowym na analizę przebiegu wydarzeń występujących w aplikacji, natomiast jest niezwykle czasochłonne i podatne na błędy. Istnieją jednak narzędzia realizujące odmienne funkcje, które częściowo pokrywają się z obszarem działania Magpie Monitor.

Przykładem takiego narzędzia może być platforma Science Logic [22], która wykorzystuje logi emitowane przez aplikacje do przeprowadzania analizy przyczyn źródłowych (ang. root cause analysis) incydentów, w których aplikacja działała nieprawidłowo. Istotną różnicą jest jednak moment, w którym owa analiza jest dokonywana. Science Logic dokonuje analizy po wystąpieniu incydentu, w którym aplikacja działała nieprawidłowo, natomiast Magpie Monitor stara się zapobiec możliwemu incydentowi, ostrzegając użytkownika przed potencjalnym zagrożeniem.

Kolejnym narzędziem, którego działanie częściowo pokrywa się z funkcjami Magpie Monitor jest Elasticsearch [23], a konkretniej Elastic Cloud [24], czyli platforma udostępniająca zarządzaną przez firmę Elastic bazę danych Elasticsearch. Platforma ta posiada mechanizm o nazwie Watcher, pozwalający na wstrzyknięcie dowolnego skryptu, który zostanie uruchomiony w momencie pojawienia się nowych logów w bazie, lub na życzenie użytkownika. Wstrzykiwany skrypt może przykładowo dokonywać analizy logów przy pomocy dużego modelu językowego, a następnie przedstawiać użytkownikowi potencjalne zagrożenia oraz rekomendacje. Tematyka ta, została poruszona w poście [25] napisanym przez Elastic Labs, czyli podfirmę należącą do Elastic. Narzędzie Watcher pozwala na osiągnięcie podobnego rezultatu, natomiast wymaga ono korzystania z chmury Elastic, a także skazuje użytkownika na niewygodny interfejs odczytu rekomendacji, który utrudnia wiele funkcjonalności, w szczególności takich jak wysyłanie powiadomień. Mechanizm watcher może służyć do szybkiego i wygodnego prototypowania rozwiązań do analizy semantycznej logów, natomiast jego rozbudowa oraz skalowanie jest ograniczone.

Pewną konkurencją dla Magpie Monitor mogą być również modele pochodzące z GPT Store firmy OpenAI, który pozwala użytkownikom na dotrenowanie modelu GPT na dodatkowych danych, dzięki czemu model lepiej porusza się w obszarze wybranego zagadnienia. Przykładem takiego modelu może być LogAnalyzer [26], którego wytrenowano dodatkowo w obszarze analizy logów. Uzyskiwane przez niego rezultaty są jednak zbliżone do domyślnego modelu GPT-4, z którego korzysta Magpie Monitor. Warto jednak podkreślić, że narzędzia podobnego typu mają spory potencjał, który uwidocznić mogą jedynie starannie dobrane dane treningowe.

Warto również wspomnieć, że dziedzina analizy semantycznej logów przy pomocy dużych modeli językowych jest lukratywna dla naukowców, co pokazuje spora liczba artykułów naukowych, które badają ten obszar. Niektóre z artykułów np. [27], proponują nawet standaryzowane testy wydajnościowe, które pomogą w ewaluacji przyszłych rozwiązań. Zainteresowanie naukowców wskazuje na fakt, iż analiza logów z użyciem dużych modeli językowych jest nietrywialna, co jest wspaniałą okazją do wytworzenia przewagi konkurencyjnej względem innych, bardziej tradycyjnych rozwiązań.

Podsumowując, dziedzina analizy logów z użyciem dużych modeli językowych jest stosunkowo nowym obszarem, którego komercyjne pokrycie jest niewielkie. Obszar ten wydaje się jednak być obiecującym, co m.in. potwierdza zainteresowanie wśród naukowców oraz sporych firm. Osobiście wierzymy, że przyszłość dąży do automatyzacji poprzez AI, w szczególności w dziedzinie analizy logów oraz detekcji anomalii, gdzie wolumen danych niejednokrotnie przekracza ludzkie możliwości analizy. Jednocześnie sądzimy, że zbierane dane często są marnowane, ponieważ nie wyciąga się z nich wniosków, które mogą być wysoce pomocne. Wysoki potencjał pokrywanego przez Magpie Monitor obszaru rynkowego oraz jego stosunkowo niskie zaspokojenie przez aktualne narzędzia sprawia, że projekt oraz uzyskany przez nas rezultat ma realny sens w kontekście dziedziny, w której jest osadzony.

# 5. Założenia wstępne {#założenia-wstępne}

## 5.1 Przeznaczenie {#przeznaczenie}

Celem systemu Magpie Monitor jest wsparcie administratorów IT w monitorowaniu usług webowych, które muszą być dostępne przez całą dobę. Każda awaria może generować poważne straty finansowe dla firmy, dlatego system został zaprojektowany, aby pomagać w zapobieganiu potencjalnym przestojom oraz przyspieszać proces monitorowania aplikacji. Rozwiązanie to nie tylko identyfikuje problemy, ale również proaktywnie wspiera działania prewencyjne, minimalizując ryzyko awarii.

Warto jednak zaznaczyć, że wykorzystanie zewnętrznego modelu, w celu generowania raportów w języku naturalnym, wiąże się z określonymi kosztami operacyjnymi. Potencjalni klienci muszą więc przeanalizować, czy koszty wdrożenia i utrzymania Magpie Monitora będą niższe od strat wynikających z ewentualnych przestojów.

Większość konkurencyjnych rozwiązań skupia się na tzw. podejściu „slice and dice”. Polega ono na
podzieleniu danych na mniejsze fragmenty i prezentowaniu ich w sposób umożliwiający analizę trendów w zbiorze danych. Magpie Monitor działa w sposób bardziej zautomatyzowany, pokazując nie tylko trendy w danych, ale także ich konsekwencje oraz potencjalne sposoby rozwiązania problemów. Jest to nowatorskie rozwiązaniem, które wyróżnia się w segmencie narzędzi do monitorowania usług webowych.

## 5.2 Nazwa i logo projektu {#nazwa-i-logo-projektu}

Projekt nosi nazwę _Magpie Monitor_, a oba człony tej nazwy zostały dobrane nieprzypadkowo. Słowo _Monitor_ jednoznacznie wskazuje na podstawową funkcjonalność systemu, jaką jest monitorowanie aplikacji. Z kolei _Magpie_ (sroka) nawiązuje do charakterystycznego zachowania tego ptaka, który zbiera różnorodne błyskotki – analogicznie do zaprojektowanego systemu, który gromadzi logi. Motyw sroki pojawia się również w logo projektu, nadając mu spójności wizualnej z nazwą:

<br>
<br>
<br>

<figure>
    <img src="/logo.png">
    <figcaption>Logotyp systemu Magpie Monitor [źródło opracowanie własne]</figcaption>
</figure>

Logo jest minimalistyczne i przedstawia czarno-białą srokę na szarym tle. Subtelnie zaokrąglone krawędzie dodają obrazowi delikatności, a poziome czarne linie w tle podkreślają czujną, wypatrującą pozę ptaka. Taka forma graficzna odzwierciedla główną ideę systemu: uważne monitorowanie oraz zbieranie danych.

Głównym celem logotypu jest ścisłe powiązanie go z nazwą systemu. Użycie maskotki w postaci sroki skutecznie realizuje ten zamysł, tworząc czytelne skojarzenie wizualne i jednocześnie dodając projektowi charakteru.

## 5.3 Analiza funkcjonalna {#analiza-funkcjonalna}

System będzie oferował funkcjonalność uwierzytelniania w celu identyfikacji użytkowników i weryfikacji, czy osoby próbujące uzyskać dostęp do systemu posiadają odpowiednie uprawnienia. Funkcja ta jest kluczowa dla zapewnienia bezpieczeństwa działania aplikacji.  
Po zalogowaniu użytkownik powinien ocenić interfejs aplikacji jako intuicyjny i estetyczny. Z tego względu raporty generowane przez system nie mogą być przedstawiane w formie surowego tekstu. Informacje w raportach zostaną podzielone na czytelne sekcje, a niektóre dane dotyczące przeanalizowanych logów będą prezentowane w formie wykresów, co zwiększy ich przejrzystość.  
Zakres przedstawionych danych w raporcie będzie zależny od opcji wybranych przez użytkownika. Przewidywane są następujące parametry konfiguracji:

- **Dokładność** – określa, jaka część logów zostanie odfiltrowana przed ich przetworzeniem przez duży model językowy.
- **Okres** – definiuje przedział czasowy, w którym wygenerowane logi mają być uwzględnione w raporcie.
- **Źródła logów** – pozwala na wybór liczby węzłów Kubernetesa oraz aplikacji, które mają być uwzględnione w procesie generowania raportu.
- **Indywidualne instrukcje wejściowe** (z ang. custom prompt) – pozwala użytkownikowi dodać własne instrukcję dla modelu, analizującego logi.

Te parametry bezpośrednio wpłyną na jakość raportu oraz związane z nim koszty. Aby użytkownicy nie musieli za każdym razem konfigurować i generować raportów, system będzie oferował funkcjonalność raportów cyklicznych, które będą tworzone automatycznie w określonych przedziałach czasowych.  
Ze względu na to, że czas generacji raportu nie jest możliwy do precyzyjnego określenia z góry, użytkownik będzie mógł przypisać kanał powiadomień (np. e-mail), na który system wyśle informację o zmianie statusu generacji raportu. Dzięki temu użytkownik zostanie poinformowany na bieżąco o zakończeniu procesu lub wystąpieniu ewentualnych problemów.

## 5.4 Przyjęte ograniczenia {#przyjęte-ograniczenia}

W realizowanym projekcie przyjęto następujące ograniczenia:

- System będzie wykorzystywał wyłącznie modeli od OpenAI, które wspierają funkcję ustrukturyzowane wyjście (ang. structured outputs). Nie przewiduje się wsparcia dla alternatywnych modeli językowych.
- Jedynym wspieranym orkiestratorem z jakiego możliwe jest zbieranie logów to Kubernetes.
- System ogranicza dostępne kanały powiadomień do e-maila, Slacka oraz Discorda. Inne formy powiadomień nie są obsługiwane w bieżącej wersji.
- System logowania zostanie zaimplementowany w oparciu o logowanie jednokrotne (ang. Single Sign-On, SSO) firmy Google. W obecnym etapie projektu nie przewiduje się wsparcia dla alternatywnych metod logowania.

# 6. Specyfikacja i analiza wymagań na produkt programowy {#specyfikacja-i-analiza-wymagań-na-produkt-programowy}

## 6.1 Użytkownicy systemu {#użytkownicy-systemu}

| Nazwa      | Opis                                                                     | Zakres funkcjonalności                                                                                                             |
| :--------- | :----------------------------------------------------------------------- | :--------------------------------------------------------------------------------------------------------------------------------- |
| Gość       | Osoba niezalogowana w systemie.                                          | Logowanie się w systemie.                                                                                                          |
| Użytkownik | Administrator lub osoba odpowiedzialna za monitorowanie stanu aplikacji. | Wylogowywanie się z systemu, przeglądanie raportów, generowanie raportów, konfigurowanie raportów, ustawianie kanałów powiadomień. |

## 6.2 Wymagania niefunkcjonalne {#wymagania-niefunkcjonalne}

Wymagania funkcjonalne zostały podzielone według modelu FURPS. Pozwala on skategoryzować wszystkie wymagania według pięciu klas (Funkcjonalności, Użyteczności, Niezawodności, Wydajności, Wsparcia).

**1. Funkcjonalność**

- System musi uniemożliwiać dostęp nieautoryzowanym użytkownikom, aby zapewnić bezpieczeństwo danych klientów.

**2. Użyteczność**

- System powinien oferować intuicyjny i wygodny interfejs użytkownika.
  Pozytywne doświadczenia użytkowników z korzystania z Magpie Monitora mają na celu zwiększenie zaangażowania i zachęcenie do dalszego korzystania z systemu.

- System musi być dostępny na urządzeniach mobilnych, tabletach i komputerach stacjonarnych, zapewniając pełną responsywność interfejsu użytkownika (UI).
  Chociaż przewiduje się, że Magpie Monitor będzie głównie używany na komputerach stacjonarnych, system powinien być równie funkcjonalny na innych urządzeniach, aby poszerzyć grono potencjalnych użytkowników.

**3. Niezawodność**

- System musi być dostępny dla użytkowników przez co najmniej 99% czasu w ciągu miesiąca, co zapewni wysoką niezawodność usług.

- Magpie Monitor musi być bardziej niezawodny niż aplikacje, które monitoruje, aby jego wyniki były wiarygodne i wartościowe.

- System nie może utracić więcej zebranych logów niż 0.5% w skali miesiąca.
  Utracenie większej liczby logów będzie skutkować raportami, które nie opisują rzeczywistego stanu monitorowanej aplikacji.

- System po awarii musi wznowić generacje raportu
  Każdy zaplanowany raport musi zostać wykonany, jeśli Magpie Monitor oraz używany model językowy jest sprawny.

**4. Wydajność**

- System powinien nie mieć ograniczeń w liczbie analizowanych logów.
  Jeśli liczba logów przekracza ograniczenia wybranego modelu, system powinien podzielić je na mniejsze partie logów, które będą wysłane osobno do modelu.

**5. Wsparcie**

- System powinien wspierać popularne przeglądarki internetowe takie jak: Google Chrome, Safari, Edge, Firefox.

- System powinien być w stanie obserwować dowolną aplikację, która została umieszczona na klastrze Kubernetesa.

## 6.3 Wymagania funkcjonalne {#wymagania-funkcjonalne}

Wymagania funkcjonalne to szczegółowe opisy tego, co system ma robić, aby spełnić potrzeby użytkownika i osiągnąć zamierzone cele.  
 Poniższa tabela prezentuje wszystkie wymagania funkcjonalne dla produktu Magpie Monitor:

| Id    | Nazwa                                                            | Opis                                                                                                                                                                                                                                                 |
| :---- | :--------------------------------------------------------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| REQ01 | Logowanie                                                        | Użytkownik powinien mieć możliwość zalogowania się do aplikacji, aby uzyskać dostęp do funkcjonalności systemu.                                                                                                                                      |
| REQ02 | Wylogowanie                                                      | Zalogowany użytkownik powinien mieć możliwość wylogowania się z aplikacji.                                                                                                                                                                           |
| REQ03 | Przeglądanie incydentów                                          | Użytkownik powinien mieć możliwość wyświetlenia listy incydentów wykrytych w wygenerowanym raporcie.                                                                                                                                                 |
| REQ04 | Przeglądanie statystyk działania systemu                         | System powinien prezentować w raportach statystyki dotyczące liczby analizowanych aplikacji, hostów, logów o różnej pilności oraz liczby odfiltrowanych logów.                                                                                       |
| REQ05 | Konfiguracja dokładności raportu                                 | System powinien umożliwiać wybór poziomu dokładności filtrowania logów, co wpływa na liczbę logów uwzględnionych w raporcie.                                                                                                                         |
| REQ06 | Konfiguracja indywidualnej dokładności dla aplikacji oraz hostów | System powinien umożliwić wybór różnych dokładności dla różnych aplikacji oraz hostów.                                                                                                                                                               |
| REQ07 | Planowanie raportów                                              | Użytkownik powinien mieć możliwość planowania cyklicznego generowania raportów w określonych przez siebie odstępach czasu.                                                                                                                           |
| REQ08 | Generacja raportu na żądanie                                     | Użytkownik powinien mieć możliwość wygenerowania raportu z logów zebranych w wybranym przedziale czasowym.                                                                                                                                           |
| REQ09 | Indywidualne dostosowanie interpretacji logów                    | Użytkownik powinien mieć możliwość personalizacji działania dużego modelu językowego, poprzez dodanie własnych instrukcji wejściowych (ang. _custom prompts_) dla każdej wybranej aplikacji oraz hosta, które będą podstawą do generowania raportów. |
| REQ10 | Wybór analizowanych hostów                                       | Użytkownik powinien mieć możliwość wskazania hostów, które mają zostać przeanalizowane w raporcie.                                                                                                                                                   |
| REQ11 | Wybór analizowanych aplikacji                                    | Użytkownik powinien mieć możliwość wskazania aplikacji, które mają zostać przeanalizowane w raporcie.                                                                                                                                                |
| REQ12 | Przypisanie powiadomień do raportu                               | System powinien umożliwiać przypisanie kanału powiadomień (Slack, Discord lub email) do raportu, aby informować użytkownika o zakończeniu generacji raportu.                                                                                         |
| REQ13 | Konfiguracja powiadomień                                         | System powinien umożliwiać dodanie, wyświetlenie, zmodyfikowanie oraz usunięcie wybranego kanału powiadomień.                                                                                                                                        |
| REQ14 | Testowanie powiadomień                                           | Użytkownik powinien mieć możliwość przetestowania kanału powiadomień w celu sprawdzenia czy działa on poprawnie.                                                                                                                                     |
| REQ15 | Wyświetlenie nazwy aplikacji, w której wykryto incydent          | Użytkownik powinien mieć możliwość wyświetlenia nazwy aplikacji, w której logach wykryto przeglądany incydent.                                                                                                                                       |
| REQ16 | Wyświetlenie nazwy hosta, w którym wykryto incydent              | Użytkownik powinien mieć możliwość wyświetlenia nazwy hosta, w którego logach wykryto przeglądany incydent.                                                                                                                                          |
| REQ17 | Rekomendowanie naprawy incydentu                                 | System powinien dostarczać rekomendacje dotyczące naprawy incydentu na podstawie wygenerowanych raportów.                                                                                                                                            |
| REQ18 | Kategoryzowanie incydentów                                       | System powinien przypisywać każdemu incydentowi odpowiednią pilność.                                                                                                                                                                                 |
| REQ19 | Streszczanie incydentu                                           | System powinien przygotowywać zwięzły opis każdego incydentu, aby dostarczyć użytkownikowi najważniejsze informacje w języku naturalnym.                                                                                                             |
| REQ20 | Nazywanie incydentów                                             | System powinien nadawać incydentom nazwy odzwierciedlające charakter problemu, jaki sugerują analizowane logi.                                                                                                                                       |
| REQ21 | Wyznaczanie zakresu czasowego incydentu                          | System powinien precyzyjnie określać przedział czasowy incydentu, uwzględniając okres, w którym zarejestrowano logi wskazujące na jego wystąpienie.                                                                                                  |
| REQ22 | Wyświetlanie źródeł incydentów                                   | Użytkownik powinien mieć możliwość wyświetlenia logów na podstawie, których został wykryty incydent.                                                                                                                                                 |
| REQ23 | Wyświetlanie listy generowanych raportów                         | Użytkownik powinien mieć możliwość wyświetlenia listy aktualnie generowanych raportów.                                                                                                                                                               |
| REQ24 | Wyświetlenie listy wygenerowanych raportów na żądanie            | Użytkownik powinien mieć możliwość przeglądania listy raportów wygenerowanych na żądanie w celu analizy historycznych danych.                                                                                                                        |
| REQ25 | Wyświetlenie listy wygenerowanych raportów cyklicznych           | Użytkownik powinien mieć możliwość przeglądania listy cyklicznie generowanych raportów, aby analizować zmiany w działaniu systemu.                                                                                                                   |
| REQ26 | Wyświetlenie listy podłączonych klastrów                         | Użytkownik powinien mieć możliwość wyświetlenia listy klastrów, z których mogą być generowane raporty.                                                                                                                                               |
| REQ27 | Oznaczenie pilności incydentu                                    | System powinien oznaczać poziom pilności każdego incydentu, aby umożliwić użytkownikowi łatwe określenie jego istotności i priorytetyzację działań.                                                                                                  |

## 6.4 Historyjki użytkownika {#historyjki-użytkownika}

W celu wskazania wymaganych funkcjonalności przez poszczególnych użytkowników systemu, zapisano historyjki użytkowników, które mają zapewnić odpowiedni kontekst i uzasadnienie poszczególnych elementów Magpie Monitora:

1. **H1 – Logowanie do systemu**  
   Jako gość,  
   Chcę zalogować się do systemu,  
   Żeby uzyskać dostęp do funkcjonalności aplikacji.

2. **H2 – Wylogowanie z systemu**  
   Jako użytkownik,  
   Chcę wylogować się z systemu,  
   Żeby zakończyć swoją sesję i zabezpieczyć dane.

3. **H3 – Przeglądanie listy incydentów**  
   Jako użytkownik,  
   Chcę przeglądać listę incydentów wykrytych w raporcie,  
   Żeby zidentyfikować potencjalne problemy w monitorowanym systemie.

4. **H4 – Wyświetlenie liczby analizowanych aplikacji**  
   Jako użytkownik,  
   Chcę zobaczyć liczbę przeanalizowanych aplikacji w raporcie,  
   Żeby wiedzieć ile aplikacji zostało przeanalizowanych.

5. **H5 – Wyświetlenie liczby analizowanych hostów**  
   Jako użytkownik,  
   Chcę zobaczyć liczbę przeanalizowanych hostów w raporcie,  
   Żeby wiedzieć ile hostów zostało przeanalizowanych.

6. **H6 – Wyświetlenie liczby incydentów**  
   Jako użytkownik,  
   Chcę zobaczyć liczbę incydentów z podziałem według ich pilności w raporcie,  
   Żeby wiedzieć ile oraz jak ważnych incydentów system odnalazł.

7. **H7 – Wyświetlenie liczby przeanalizowanych logów**  
   Jako użytkownik,  
   Chcę zobaczyć liczbę przeanalizowanych logów z aplikacji oraz hostów w raporcie,  
   Żeby określić ilość danych przetworzonych przez model językowy.

8. **H8 – Przedstawienie nazwy hosta, na którym wystąpiło najwięcej incydentów**  
   Jako użytkownik,  
   Chcę zobaczyć nazwę hosta, na którym wystąpiło najwięcej incydentów w ramach raportu,  
   Żeby lepiej zrozumieć, który host powinien być szczególnie obserwowany.

9. **H9 – Konfiguracja dokładności raportu**  
   Jako użytkownik,  
   Chcę dostosowywać poziomy dokładności filtrowania logów,  
   Żeby uzyskać obserwacje o wymaganym poziomie szczegółowości.

10. **H10 – Konfiguracja dokładności aplikacji i hostów**  
    Jako użytkownik,  
    Chcę dostosowywać poziomy dokładności osobno dla każdej aplikacji i dla każdego hosta uwzględnionego w raporcie,  
    Żeby mieć wpływ jakie hosty oraz aplikacje mają być przeanalizowane w większym stopniu, a jakie w mniejszym.

11. **H11 – Planowanie raportów**  
    Jako użytkownik,  
    Chcę planować cykliczne generowanie raportów, co określony przedział czasu,  
    Żeby otrzymywać regularne informacje o stanie monitorowanego systemu.

12. **H12 – Generacja raportu na żądanie**  
    Jako użytkownik,  
    Chcę generować raport z wybranego przedziału czasu,  
    Żeby szybko przeanalizować zebrane w tym przedziale dane.

13. **H13 – Personalizacja interpretacji logów**  
    Jako użytkownik,  
    Chcę dodawać własne instrukcje dla modelu językowego,  
    Żeby raporty były dostosowane do specyfiki moich aplikacji i hostów.

14. **H14 – Wybór analizowanych hostów**  
    Jako użytkownik,  
    Chcę wskazywać hosty, które będą źródłem danych do raportu,  
    Żeby skupić się na konkretnych elementach infrastruktury.

15. **H15 – Wybór analizowanych aplikacji**  
    Jako użytkownik,  
    Chcę wskazywać aplikacje jako źródło danych,  
    Żeby raport dotyczył tylko wybranych komponentów systemu.

16. **H16 – Konfiguracja powiadomień**  
    Jako użytkownik,  
    Chcę przypisać kanał powiadomień (Slack, Discord lub email) do raportu,  
    Żeby otrzymać informację o zakończeniu jego generacji.

17. **H17 – Dodanie nowego kanału powiadomień**  
    Jako użytkownik,  
    Chcę móc dodać nowy kanał powiadomień w systemie,  
    Żeby móc przypisać go do generowanego raportu.

18. **H18 – Usuwanie kanału powiadomień**  
    Jako użytkownik,  
    Chcę móc usunąć kanał powiadomień,  
    Żeby usunąć błędnie wprowadzony kanał.

19. **H19 – Edytowanie kanału powiadomień**  
    Jako użytkownik,  
    Chcę móc edytować dodany kanale powiadomień,  
    Żeby poprawić lub zaktualizować błędnie wprowadzone dane.

20. **H20 – Testowanie powiadomień**  
    Jako użytkownik,  
    Chcę móc testować kanał powiadomień,  
    Żeby upewnić się, że działa poprawnie.

21. **H21 – Wyświetlenie nazwy aplikacji incydentu**  
    Jako użytkownik,  
    Chcę wyświetlić nazwę hosta, w którym wykryto incydent,  
    Żeby lepiej zrozumieć przyczyny problemów.

22. **H22 – Wyświetlenie nazwy hosta incydentu**  
    Jako użytkownik,  
    Chcę wyświetlić nazwę hosta, w którym wykryto incydent,  
    Żeby lepiej zrozumieć przyczyny problemów.

23. **H23 – Rekomendacje naprawy incydentu**  
    Jako użytkownik,  
    Chcę otrzymywać rekomendacje dotyczące naprawy incydentu,  
    Żeby szybko podjąć odpowiednie działania.

24. **H24 – Kategoryzacja incydentów**  
    Jako użytkownik,  
    Chcę, aby system kategoryzował incydenty według pilności,  
    Żeby łatwiej analizować rodzaje problemów.

25. **H25 – Streszczenie incydentu**  
    Jako użytkownik,  
    Chcę otrzymać zwięzły opis incydentu,  
    Żeby szybko zrozumieć jego istotę.

26. **H26 – Nazywanie incydentów**  
    Jako użytkownik,  
    Chcę, aby incydenty miały trafne nazwy odzwierciedlające ich charakter,  
    Żeby łatwiej identyfikować problemy.

27. **H27 – Wyznaczanie czasu, w jakim wystąpił incydent**  
    Jako użytkownik,  
    Chcę, aby system określał przedział czasowy wystąpienia incydentu,  
    Żeby precyzyjnie osadzić go w kontekście zdarzeń systemowych.

28. **H28 – Wyświetlanie źródeł incydentu**  
    Jako użytkownik,  
    Chcę wyświetlać logi na podstawie, których został wykryty incydent,  
    Żeby poznać kontekst związany z incydentem.

29. **H29 – Wyświetlanie listy generowanych raportów**  
    Jako użytkownik,  
    Chcę przeglądać listę raportów w trakcie generacji,  
    Żeby przejrzeć oczekujące na generacje raporty.

30. **H30 – Wyświetlanie listy raportów wygenerowanych na żądanie**  
    Jako użytkownik,  
    Chcę przeglądać raporty wygenerowane na żądanie,  
    Żeby analizować historyczne dane.

31. **H31 – Wyświetlanie listy raportów cyklicznych**  
    Jako użytkownik,  
    Chcę przeglądać listę cyklicznie generowanych raportów,  
    Żeby obserwować zmiany w czasie.

32. **H32 – Wyświetlanie listy podłączonych klastrów**  
    Jako użytkownik,  
    Chcę widzieć listę klastrów dostępnych w systemie,  
    Żeby móc wybrać te, z których generuje się raporty.

33. **H33 – Oznaczenie pilność incydentu**  
    Jako użytkownik,  
    Chcę widzieć poziom pilności każdego incydentu,  
    Żeby priorytetyzować działania naprawcze.

## 6.5 Mapowanie wymagań funkcjonalnych na historyjki użytkownika {#mapowanie-wymagań-funkcjonalnych-na-historyjki-użytkownika}

Zależność między wymaganiami funkcjonalnymi, a historyjkami użytkownika przedstawia poniższy diagram:

<figure>
 <img src=/mapping-requirements-to-user-stories.svg>
<figcaption> Diagram przedstawiający mapowanie się wymagań funkcjonalnych na historyjki użytkownika [źródło opracowanie własne]</figcaption>
</figure>

# 7. Projekt produktu programowego {#projekt-produktu-programowego}

## 7.1. Wykorzystane technologie i narzędzia {#wykorzystane-technologie-i-narzędzia}

W celu realizacji projektu zdecydowano się użyć następujące technologie:

**Kubernetes** [[1]](#ref1)– najbardziej dojrzały i powszechnie stosowany orkiestrator rozproszonych systemów opartych na kontenerach, szeroko wykorzystywany w zastosowaniach komercyjnych.

**Docker** [[2]](#ref2)– najpopularniejsze narzędzie i ekosystem do budowania oraz uruchamiania kontenerów aplikacji.

**Golang** [[3]](#ref3) – język programowania umożliwiający tworzenie szybkich, odpornych na wycieki pamięci i wielowątkowych mikroserwisów bez konieczności używania dodatkowych frameworków do budowy aplikacji webowych. Dodatkowo, ekosystem Go zapewnia skuteczną integrację z interfejsem Kubernetesa, co jest kluczowe przy zbieraniu logów z klastra komputerowego zarządzanego przez Kubernetes.

**Fx** [[4]](#ref4) – biblioteka, która dostarcza funkcjonalność wstrzykiwania zależności do aplikacji w Go. Pozwala na standaryzację struktury mikroserwisów oraz większą reużywalność kodu.

**Java** [[5]](#ref5) – dojrzały i popularny język programowania, który dzięki bogatej dokumentacji i licznej społeczności znacząco przyspiesza proces rozwoju oprogramowania.

**Spring Boot** [[6]](#ref6) – popularny framework backendowy przeznaczony do budowy aplikacji webowych w architekturze REST. Oferuje sprawdzone rozwiązania w zakresie bezpieczeństwa, routingu oraz mapowania obiektowo-relacyjnego (ORM).

**Typescript** [[7]](#ref7) – język programowania, rozwijający język JavaScript o dodatkową składnie. Wprowadzone modyfikacje pozwalają na uniknięcie błędów związanych z brakiem silnego typowania.

**React** [[8]](#ref8)- framework frontendowy użyty do stworzenia klienta aplikacji. React jest najszerzej wspieranym frameworkiem do tworzenia aplikacji w architekturze SPA. To pozwala na znacznie łatwiejsze zarządzanie stanem aplikacji.

**Sass** [[9]](#ref9) - rozszerzenie klasycznego CSS, które ubogaca podstawową składnie o funkcjonalności minimalizujące duplikacje kodu, poprawiając przy tym czytelność pliku.

**Vite** [[10]](#ref10) – nowoczesny narzędzie do budowania frontendowych aplikacji webowych, które oferuje szybkie ładowanie modułów podczas rozwoju oraz efektywne budowanie w środowiskach produkcyjnych.

**PostgreSQL** [[11]](#ref11) - relacyjna baza danych, która została użyta do przechowywania informacji związanych z ustawieniami użytkownika oraz informacji o monitorowanym systemie, które cechują się możliwości ich normalizacji.

**MongoDB** [[12]](#ref12) - dokumentowa baza danych, w której zostaną przechowywane wygenerowane raporty w języku naturalnym. Raporty takie są długimi dokumentami, które nie wymagają spójności w każdym momencie oraz które ciężko byłoby efektywnie przechowywać i przetwarzać w niedokumentowej bazie danych.

**Kafka** [[13]](#ref13) - jedna z najpopularniejszych platform do strumieniowego przetwarzania danych i kolejkowania zdarzeń. Jej zastosowanie pozwala na uniezależnienie działania mikroserwisów od siebie, zapewniając efektywną komunikację pomiędzy nimi.

**ElasticSearch** [[14]](#ref14) - czyli nierelacyjna, łatwo skalowalna baza danych, która stała się biznesowym standardem do przechowywania logów.

**Redis** [[15]](#ref15) – szybka, nierelacyjna baza danych typu klucz-wartość, używana w projekcie jako mechanizm pamięci podręcznej, co przyspiesza dostęp do często wykorzystywanych danych oraz zmniejsza obciążenie głównych baz danych.

**Nginx** [[16]](#ref16) – reverse proxy i serwer webowy, który wspiera aplikację w obsłudze ruchu sieciowego, zwiększając jej skalowalność i wydajność.

## 7.2 Architektura systemu {#architektura-systemu}

### 7.2.1 Model C4 \- poziom 1 {#model-c4---poziom-1}

Diagram kontekstu (C1) przedstawia ogólny obraz interakcji pomiędzy kluczowymi elementami systemu i jego otoczeniem. W skład diagramu C1 dla Magpie Monitora wchodzą następujące elementy:

**Monitorowany system** \- obserwowane oprogramowanie, którego logi mają zostać przeanalizowane i na którego stabilności zależy klientowi.

**Magpie Monitor** \- główny system zapewniający mechanizm zbierania oraz analizy logów. Oferuje użytkownikowi klienta webowego, za pomocą, którego może zobaczyć uzyskane wyniki \- raporty.

**Duży model językowy** \- zewnętrzny model, który jest używany przez Magpie Monitor do generowania raportów na podstawie logów.

**Slack**, **Discord, Email**\- zewnętrzne usługi internetowe służące m.in. do komunikacji między członkami danej organizacji. W systemie Magpie Monitor są wykorzystywane jako dostępne kanały powiadomień.

<figure>
    <img src="/magpie-monitor-c4-context.drawio.svg">
    <figcaption> Diagram kontekstu [źródło opracowanie własne]</figcaption>
</figure>

### 7.2.2 Model C4 \- poziom 2 {#model-c4---poziom-2}

<figure>
    <img src="/container-diagram-latest.svg">
    <figcaption> Diagram kontenerów [źródło opracowanie własne]</figcaption>
</figure>

Poziom 2 w modelu C4 oznacza diagram kontenerów i przedstawia on wszystkie mikroserwisy, ich zewnętrze zależności oraz kanały komunikacji.

System został zaprojektowany wokół wydarzeń, w związku z tym mikroserwisy są zorientowane na wydarzenia/procesy. W ramach tego podziału wydzielone zostały następujące serwisy:

Logs Ingestion service: Odpowiada za agregowanie logów dostarczonych przez agenta zbierającego je z klastra klienta.

Reports service: Odpowiada za generowanie raportów na podstawie logów zebranych przez Logs Ingestion service oraz konfiguracji dostarczonych przez Management Service.

Metadata service: Odpowiedzialny za zbieranie i przetwarzanie bieżącego stanu klastra (aplikacji i hostów będących częścią klastra) na podstawie danych dostarczanych przez agenta zainstalowanego w systemie klienta.

Management Service: Odpowiedzialny za uwierzytelnianie użytkownika, wysyłanie powiadomień, konfigurację raportów (wykorzystująć stan klastra z Metadata service) oraz cykliczne żądanie generowania raportów od Reports service. Serwis ten również pełni rolę “Backend For Frontend” (BFF), który wystawie wygodne i zoptymalizowane RESTowe API dla klienta webowego.

Dodatkowo, aby nie naruszać zasad bezpieczeństwa, które klient może mieć w swoim systemie, wykorzystujemy agenta, który będąc częścią klastra klienta wysyła logi do zewnętrznej sieci, w której znajdują się mikroserwisy Magpie Monitor. Stosując takie podejście klient nie musi otwierać dodatkowych portów sieciowych w swoim systemie. Dzięki temu tylko kanał komunikacji z Logs Ingestion Service (w tym przypadku Apache Kafka) musi mieć globalnie routowalny adres (i tym samym otwarty port.)

Komunikacja pomiędzy mikroserwisami odbywa się za pomocą brokerów wiadomości (Apache Kafka). Podejście to znacząco zwiększa niezawodność i identyfikowalność (traceability) w systemie. Aby dostarczać większość funkcjonalności w dowolnym momencie wystarczy aby działał wyłącznie Management Service, inne serwisy wykonają jego żądania w momencie zakończenia ich awarii, ale same wiadomości nie zostaną utracone lub nie spowodują awarii samego Management Service.

Identyfikowalność w ramach komunikacji między serwisami została zrealizowana z wykorzystaniem identyfikatora korelacji (correlation id), który służy do śledzenia całego procesu, w ramach którego przesyłane są wiadomości z tym samym identyfikatorem pomiędzy wieloma serwisami. Dzięki temu serwisy są w stanie rozpoznać które żądanie zostało zaktualizowane, lub jaki jest stan danego procesu.

Ze względu na zorientowanie serwisów na wydarzenia, każdy z nich musi mieć bezpośredni dostęp do wszystkich danych wymaganych do realizacji zadań danego serwisu. Oznacza to, że każdy z nich ma swoje bazy danych, które są tworzone oraz aktualizowane na podstawie wydarzeń, które dany serwis otrzymał.

Jednocześnie, serwis w jakim utworzony został dany rekord/informacja oryginalnie jest odpowiedzialny za zachowanie spójności nadając unikalne identyfikatory, które są powielane w innych serwisach wykorzystujących te dane.

### 7.2.3 Diagram wdrożenia {#diagram-wdrożenia}

<figure>
    <img src="/deployment-diagram-latest.svg">
    <figcaption> Diagram wdrożenia [źródło opracowanie własne]</figcaption>
</figure>

Wdrożenie systemu zakłada dwa podsystemy: “Magpie Monitor Cloud”, który oznacza infrastrukturę, na której wdrażane są wszystkie mikroserwisy oraz serwer klienta webowego. Drugim podsystemem jest system klienta, w którym musi zostać zainstalowany **agent** zbierający logi z klastra Kubernetesa będącego częścią jego infrastruktury.

W ramach systemu “Mapgie Monitor Cloud” mikroserwisy oraz ich bazy danych zostały wydzielone do osobnych hostów aby zminimalizować ryzyko awarii wielu serwisów jednocześnie. Serwisy takie jak **Logs Ingestion service** oraz **Cluster Metadata service** nieprzerwanie przetwarzają duże ilości danych, w związku z czym ich obciążenie jest niezależne od obciążenia innych serwisów, w których obciążenie jest związane z obecnym ruchem sieciowym.

Wiadomości dostarczane do **Logs ingestion service** i **Cluster metadata service** również są realizowane poprzez zewnętrznego brokera, który pozwala na zmniejszenie obciążenia brokera odpowiedzialnego za wewnętrzną komunikację oraz zastosowanie osobnych reguł bezpieczeństwa, tak, aby mógł być on bezpiecznie wystawiona do internetu (aby agent mógł się z nim komunikować)

Tak jak wspomniano przy Diagramie C2, prawie każdy z serwisów ma swoją bazę, którą populuje i aktualizuje na podstawie własnych działań oraz wydarzeń otrzymywanych od innych serwisów. Wyjątkiem jest baza danych odpowiedzialna za przechowywanie logów.

Przez duży wolumen danych nieopłacalnym było przesyłanie ich za pomocą brokera wiadomości, więc zdecydowano się wykorzystać mechanizm dostarczany przez ElasticSearch. W tym przypadku ElasticSearch umożliwia stworzenie klastra swoich instancji, poprzez dodanie replik oraz jednej instancji do której można zapisywać nowe dane. Rozwiązanie to jest idealne w przypadku Magpie Monitor, ponieważ dwoma serwisami, które wykorzystują bazę logów są **Logs Ingestion service** oraz **Reports service.** Gdzie **Logs Ingestion service** jedynie zapisuje logi do bazy, a **Reports service** jedynie je odczytuje (bez ich modyfikacji). Takie rozwiązanie oferuje separacje odpowiedzialności utrzymywania spójności przez Logs Ingestion, jednocześnie nie powodując zwiększonego obciążenia na instancji wykorzystywanej do odczytu przez **Reports service.**

Wdrożenie na klastrze klienta zakłada zainstalowania agentów skonfigurowanych odpowiednio do zbierania logów i metadanych z aplikacji oraz hostów w klastrze Kubernetesa. Aby zachować trwałość danych w przypadku chwilowej awarii, dodatkowym serwisem, który musi być zainstalowany razem z agentami jest baza danych służąca do utrzymywania metadanych z klastra, która w tym przypadku jest instancją Redisa. Redis został wybrany ze względu na szybki dostęp do danych oraz niskie zużycie zasobów, co jest kluczowe w systemie klienta.

### 7.2.4 Model C3 \- poziom 3 \- Agent {#model-c3---poziom-3---agent}

<figure>
    <img src="/agent/agent-components-transparent.png">
    <figcaption> Agent: Diagram komponentów [źródło opracowanie własne]</figcaption>
</figure>

Agent instalowany jest na klastrze Kubernetes klienta i odpowiada za zbieranie metadanych i logów aplikacji oraz hostów, które następnie przesyła do brokera Kafki. Implementacyjnie Agent składa się z dwóch modułów podrzędnych, tj. Pod Agent oraz Node Agent. Pod Agent zbiera logi oraz metadane aplikacji, wykorzystując przy tym API klastra Kubernetes. Node Agent zbiera metadane oraz logi z hostów, wykorzystując przy tym API dostępowe systemu plików Linux. Postęp czytania logów hostów dla danego pliku jest zapisywany w zewnętrznej bazie danych Redis, dzięki czemu przesyłane dane nie są powtarzane nawet w przypadku tymczasowej awarii mechanizmu zbierania danych.

### 7.2.5. Model C3 \- poziom 3 \- Logs Ingestion Service {#model-c3---poziom-3---logs-ingestion-service}

<figure>
    <img src="/logs-ingestion/logs-ingestion-components.svg">
    <figcaption> Logs Ingestion: Diagram komponentów [źródło opracowanie własne]</figcaption>
</figure>

Logs Ingestion Service zajmuje się zbieraniem i przetwarzaniem kolejno logów aplikacji i hostów w Application Logs Queue i Node Logs Queue. Logi te są spłaszczone do formy, którą można efektywnie zapisać w bazie logów wykorzystując kolejno Application Logs Repository oraz Node Logs Repository. Wymienione powyżej serwisy pełnią funkcje abstrakcji (udostępniają interfejs) dla konkretnych technologii, które były wykorzystane w ramach Magpie Monitor. Takie rozwiązanie pozwala na łatwą podmianę użytego brokera wiadomości lub bazy danych jeżeli zajdzie taka potrzeba.

### 7.2.6. Model C3 \- poziom 3 \- Report Service {#model-c3---poziom-3---report-service}

<figure>
    <img src="/reports/reports-components.svg">
    <figcaption> Reports Service: Diagram komponentów [źródło opracowanie własne]</figcaption>
</figure>

**Reports Service** pełni 4 podstawowe funkcjonalności.

1. Zbieranie zgłoszeń o wygenerowanie raportów i wysyłanie komunikatów o gotowym raporcie, lub błędzie w trakcie generowania raportu wykorzystując **ReportsHandler**.
2. **Przygotowanie logów do raportu.** W ramach każdego z raportów można sprecyzować źródła logów i ich dokładność. Aby wykonać to zadanie, serwis wykorzystuje **Node Logs Repository** oraz **Application Logs Repository** aby pobrać logi z bazy, a także Accuracy **Filter,** aby przefiltrować logi dla każdej aplikacji i hosta na podstawie sprecyzowanych konfiguracji.
3. **Stworzenie raportu z wykorzystaniem modelu językowego**. To zadanie wymaga spakowania logów w paczki, które model może jednocześnie przetworzyć. Po przesłaniu logów do modelu, serwis wykorzystuje **Batch Pollera** aby obserwować, kiedy określona paczka była przetworzona. Batch Poller przechowuje wszystkie wykonywane zadania przetwarzania paczek logów przy użyciu **Scheduled Jobs Repository.**
4. Sformułowanie raportu polega na zebraniu wszystkich rezultatów z modelu, transformacji ich do dokumentów zawierających odpowiednie metadane wykorzystując **Node Insights Generator** i **Application Insights Generator** oraz scalenie zduplikowanych incydentów z użyciem **Incident Merger**. Następnie raport jest zapisywany przy użyciu **Reports Repository**

**Reports Service** posiada wysokopoziomową logikę odpowiedzialną za generowanie raportów, ale jest niezależny od kanałów komunikacji, ponieważ za komunikację z innymi mikroserwisami odpowiedzialny jest **ReportsHandler.**

**Application Insights Generator** oraz **Node Insights Generator** stanowią interfejsy wystawiane innym serwisom i pozwalają na implementację abstrakcji nad modelem językowym użytym do generowania raportów.

### 7.2.7. Model C3 \- poziom 3 \- Metadata Service {#model-c3---poziom-3---metadata-service}

<figure>
    <img src="/metadata-service/metadata-service-components.svg">
    <figcaption> Metadata Service: Diagram komponentów [źródło opracowanie własne]</figcaption>
</figure>

Metadata Service jest systemem odpowiedzialnym za zbieranie, przechowywanie oraz przetwarzanie metadanych o klastrach, aplikacjach oraz hostach.

Metadane to informacje o aktualnie działających klastrach, o aplikacjach działających na danym klastrze, a także o hostach należących do danego klastra oraz skonfigurowanych dla nich plikach z logami.

Kluczowe funkcjonalności serwisu to:

- pobieranie metadanych z brokera Kafki, aby następnie zapisać je w dokumentowej bazie MongoDB
- wsadowe przetwarzanie zebranych danych, mające na celu wykrycie zmian w metadanych aplikacji, hostów oraz klastrów. Jeśli zmiana stanu zostanie wykryta, emitowane jest wydarzenie, którego treścią jest nowy stan. Podejście takie znacznie zmniejsza obciążenie serwisu który korzysta z metadanych, ponieważ dyskretyzuje dane ciągłe bez straty informacji, jednocześnie ograniczając liczbę próbek. Przeszukiwanie oraz scalanie danych jest wtedy znacznie szybsze niż w przypadku danych ciągłych.

### 7.3.8. Model C3 \- poziom 3 \- Management Service {#model-c3---poziom-3---management-api}

Management Service jest bramą wejściową do **Magpie Monitor Cloud**, jego zadaniem jest integracja systemu generowania raportów oraz metadanych, aby następnie wystawić przesyłane przez nie dane w formie API, z którego korzysta warstwa prezentacji zawarta w aplikacji klienckiej.

Serwis składa się z pięciu głównych komponentów, których diagramy zaprezentowano poniżej:

- podsystem raportów
- podsystem metadanych
- podsystem powiadomień
- podsystem klastrów
- podsystem uwierzytelniania, autoryzacji oraz danych o użytkowniku

<figure>
    <img src="/management-service/management-service-reports.svg">
    <figcaption> Management Service: Diagram komponentów podsystemu raportów [źródło opracowanie własne]</figcaption>
</figure>

Podsystem raportów jest odpowiedzialny za komunikacje z **Report Service,** z którym komunikuje się za pośrednictwem brokera Kafki, przez którego przesyłane są dane konfiguracyjne do generowania raportu, a następnie zwracany jest wygenerowany raport bądź odpowiedni błąd generacji, które przechowywane są w dokumentowej bazie danych MongoDB. Zapisywanie raportów w bazie sprawia, że serwis może zwracać użytkownikowi dane o raportach nawet w przypadku awarii **Report Service**.

Podsystem ten przechowuje również metadane o wykonanych generacjach raportów, które mogą służyć jako dane audytowe.

Wygenerowane raporty udostępniane są przez interfejs API, z którego może skorzystać aplikacja kliencka.

<figure>
    <img src="/management-service/management-service-metadata.svg">
    <figcaption> Management Service: Diagram komponentów podsystemu metadanych [źródło opracowanie własne]</figcaption>
</figure>

Podsystem metadanych odpowiada za odbieranie wydarzeń sygnalizujących zmianę stanu klastrów, aplikacji oraz hostów. Zagregowane dane są przechowywane w dokumentowej bazie MongoDB, a następnie są udostępniane innym podsystemom przez interfejs programistyczny.

<figure>
    <img src="/management-service/management-service-notifications.svg">
    <figcaption> Management Service: Diagram komponentów podsystemu powiadomień [źródło opracowanie własne]</figcaption>
</figure>

Podsystem powiadomień odpowiada za konfigurację odbiorników powiadomień (ang. receivers) oraz komunikację z zewnętrznymi API w celu wysyłania powiadomień. Podsystem powiadomień udostępnia interfejs programistyczny dla podsystemu raportów, pozwalając mu na przesyłanie powiadomień dot. generacji raportu zgodnie ze skonfigurowanymi przez użytkownika odbiornikami.

<figure>
    <img src="/management-service/management-service-clusters.svg">
    <figcaption> Management Service: Diagram komponentów podsystemu klastrów [źródło opracowanie własne]</figcaption>
</figure>

Podsystem klastrów odpowiada za konfigurowanie ustawień raportu oraz odbiorników powiadomień dla wybranego klastra. Podsystem ten komunikuje się również z podsystemem metadanych, pobierając metadane o klastrach, które są udostępniane w formie API.

<figure>
    <img src="/management-service/management-service-authentication.svg">
    <figcaption> Management Service: Diagram komponentów podsystemu uwierzytelniania [źródło opracowanie własne]</figcaption>
</figure>

Podsystem uwierzytelniania odpowiedzialny jest za uwierzytelnianie użytkownika przy pomocy protokołu OAuth, gdzie dostawcą uwierzytelniania jest firma Google.

### 7.2.9. Model C3 \- poziom 3 \- Client {#model-c3---poziom-3---client}

## 7.3. Bazy danych mikroserwisów {#bazy-danych-mikroserwisów}

### 7.3.1 Baza logów (Logs Ingestion Service) {#baza-logów}

**ApplicationLogsDocument**
Dokument zawierający logi i metadane aplikacji z jakiej zostały zebrane.

| Nazwa atrybutu | Znaczenie                                                                                        | Dziedzina |
| :------------- | :----------------------------------------------------------------------------------------------- | :-------- |
| id             | Unikalny identyfikator dokumentu w ramach danego indeksu                                         | string    |
| clusterId      | Identyfikator klastra, z którego pochodzi log                                                    | string    |
| kind           | Rodzaj zasobu w Kubernetesie, np. Statefull Set albo Deployment                                  | string    |
| collectedAtMs  | Czas kiedy dany log był zebrany w milisekundach od początku epoki                                | int64     |
| namespace      | Przestrzeń nazw w jakiej działała aplikacja                                                      | string    |
| podName        | Nazwa poda, z którego pochodzi dany log                                                          | string    |
| containerName  | Nazwa kontenera, z którego pochodzi dany log                                                     | string    |
| image          | Identyfikator obrazu na podstawie, którego był uruchomiony kontener, z którego pochodzi dany log | string    |
| content        | Treść loga                                                                                       | string    |

**NodeLogsDocument**
Dokument zawierający logi i metadane hosta z jakiego zostały zebrane.

| Nazwa atrybutu | Znaczenie                                                         | Dziedzina |
| :------------- | :---------------------------------------------------------------- | :-------- |
| id             | Unikalny identyfikator dokumentu w ramach danego indeksu          | string    |
| clusterId      | Identyfikator klastra, z którego pochodzi log                     | string    |
| kind           | Rodzaj zasobu w Kubernetesie, np. Node                            | string    |
| collectedAtMs  | Czas kiedy dany log był zebrany w milisekundach od początku epoki | int64     |
| name           | Nazwa hosta (na podstawie nazwy systemowej)                       | string    |
| filename       | Nazwa pliku z jakiego pochodzi log                                | string    |
| content        | Treść loga                                                        | string    |

<figure>
    <img src="/logs-ingestion/logs-ingestion-db.svg">
    <figcaption> Logs Ingestion: Diagram fizyczny bazy logów [źródło opracowanie własne]</figcaption>
</figure>

Logi zebrane z klastra klienta są przechowywane w nierelacyjnej bazie danych \- ElasticSearch. W ramach tej bazy danych “kolekcje” podzielone są na indeksy, i to w ramach tych indeksów przeszukiwane są dane podczas zapytań. W związku z tym kluczowe jest stworzenie indeksów, które oddają typowe zapytania aplikacji.

Typowym zapytaniem do bazy jest pobieranie logów do raportu. Parametrami tego zapytania jest klaster oraz przedział czasowy.

Na podstawie tej charakterystyki, indeksy zostały podzielone ze względu na klaster, miesiąc w jakim zebrano dany log oraz rodzaj logów (aplikacji lub hosta).

Indeksy przyjmują postać {nazwa-klastra}-applications-{mm.rrrr} dla logów z aplikacji oraz {nazwa-klastra}-nodes-{mm.rrrr} dla logów z hostów.

W kolekcji z logami w hostów przechowywane są dokumenty typu **NodeLogsDocument**, natomiast logi aplikacji są przechowywane w postaci dokumentu **ApplicationLogsDocument.**

Dzięki temu rozwiązaniu minimalizujemy liczbę przeglądanych logów podczas generowania raportów, tym samym zmniejszając obciążenie system i czas realizacji takiego zapytania, co jest kluczowe ze względu na objętość danych.

### 7.3.2 Baza raportów (reports service) {#baza-danych-raportów}

#### Report

Obiekt ten przechowuje raport od momentu zażądania wygenerowania do momentu wygenerowania przez model językowy.
Jest tworzony w momencie otrzymania żądania od Management Service w którym znajduje się correlationId, na podstawie którego można go powiązać z samym żądaniem.

| Nazwa atrybutu                         | Znaczenie                                                                                               | Dziedzina                                 |
| :------------------------------------- | :------------------------------------------------------------------------------------------------------ | :---------------------------------------- |
| id                                     | Unikalny identyfikator raportu                                                                          | string                                    |
| correlationId                          | Identyfikator korelacji (z żądania na podstawie, którego został wygenerowany raport)                    | string                                    |
| status                                 | Aktualny status raportu.                                                                                | string                                    |
| sinceMs                                | Data określająca początek okresu analizowanego przez raport. Wyrażona w milisekundach od początku epoki | number                                    |
| toMs                                   | Data określająca koniec okresu analizowanego przez raport. Wyrażona w milisekundach od początku epoki   | number                                    |
| requestedAt                            | Data zażądania raportu. Wyrażona w milisekundach od początku epoki                                      | number                                    |
| scheduledGenerationAtMs                | Data wygenerowania raportu. Wyrażona w milisekundach od początku epoki                                  | number                                    |
| title                                  | Tytuł raportu.                                                                                          | string                                    |
| nodeReports                            | Incydenty z hostów pogrupowane po hostcie i konfiguracji z jakimi został wygenerowany.                  | NodeReport\[\]                            |
| applicationReports                     | Incydenty z aplikacji pogrupowane po aplikacji i konfiguracji z jakimi został wygenerowany.             | ApplicationReport\[\]                     |
| totalApplicationEntries                | Liczba logów z aplikacji, uwzględnionych podczas generowani raportu                                     | number                                    |
| totalNodeEntries                       | Liczba logów z hostów, uwzględnionych podczas generowani raportu                                        | number                                    |
| urgency                                | Poziom “krytyczności” całego raportu                                                                    | string                                    |
| scheduledApplicationInsights           | Zaplanowane do przetworzania (w ramach raportu) logi i konfiguracje aplikacji                           | ScheduledApplicationInsights\[\]          |
| scheduledNodeInsights                  | Zaplanowane do przetworzania (w ramach raportu) logi i konfiguracje hostów                              | ScheduledNodeInsights\[\]                 |
| analyzedApplications                   | Liczba aplikacji, przeanalizowanych w ramach raportu                                                    | number                                    |
| analyzedNodes                          | Liczba hostów, przeanalizowanych w ramach raportu                                                       | number                                    |
| scheduledApplicationIncidentMergerJobs | Zaplanowane zadania związane ze scalaniem zduplikowanych incydentów aplikacji                           | ScheduledApplicationIncidentMergerJob\[\] |
| scheduledNodeIncidentMergerJobs        | Zaplanowane zadania związane ze scalaniem zduplikowanych incydentów hostów                              | ScheduledNodeIncidentMergerJob\[\]        |

#### ApplicationReport

Przechowuje incydenty występujące dla konkretnej aplikacji w ramach **Report**. Pozwala na połączenie informacji dotyczącej konfiguracji aplikacji przy generowaniu raportu wraz z incydentami z danej aplikacji.

| Nazwa atrybutu  | Znaczenie                                                                                              | Dziedzina                |
| :-------------- | :----------------------------------------------------------------------------------------------------- | :----------------------- |
| applicationName | Nazwa aplikacji                                                                                        | string                   |
| accuracy        | Dokładność raportu aplikacji                                                                           | string                   |
| customPrompt    | Własne dodatkowe wejście do modelu językowego podczas interpretacji logów z aplikacji w ramach raportu | string                   |
| incidents       | Lista incidentów z aplikacji                                                                           | ApplicationIncidents\[\] |

#### ApplicationIncident

Incydent aplikacji wraz z metadanymi pozwalającymi na stwierdzenie na podstawie których logów dany incydent został wykryty. Dokument ten został zdenormalizowany (posiada źródła incydentu (**ApplicationIncidentSources**) oraz nazwę aplikacji i konfigurację aplikacji, ponieważ pobieranie pojedyneczego incydentu wraz z jego metadanymi jest typowym zapytaniem, a denormalizacja pozwala na uniknięcie częstej operacji złączania dokumentów.

| Nazwa atrybutu  | Znaczenie                                                                                              | Dziedzina                     |
| :-------------- | :----------------------------------------------------------------------------------------------------- | :---------------------------- |
| id              | Unikalny identyfikator incydentu                                                                       | string                        |
| title           | Tytuł incydentu                                                                                        | string                        |
| customPrompt    | Własne dodatkowe wejście do modelu językowego podczas interpretacji logów z aplikacji w ramach raportu | string                        |
| clusterId       | Identyfikator klastra na którym wystąpił dany incydent                                                 | string                        |
| applicationName | Nazwa aplikacji dla której wystąpił dany incydent                                                      | string                        |
| category        | Kategoria incydentu                                                                                    | string                        |
| summary         | Podsumowanie incydentu                                                                                 | string                        |
| recommendation  | Rekomendacja odnośnie rozwiązania incydentu                                                            | string                        |
| urgency         | Krytyczność incydentu                                                                                  | string                        |
| sources         | Źródła na podstawie, których dany incydent był wykryty                                                 | ApplicationIncidentSource\[\] |

#### ScheduledApplicationInsights

Reprezentuje wykonywane zadanie generowania **obserwacji** z aplikacji.

| Nazwa atrybutu           | Znaczenie                                                                                               | Dziedzina                           |
| :----------------------- | :------------------------------------------------------------------------------------------------------ | :---------------------------------- |
| scheduledJobIds          | Identyfikatory zadań związanych z wykrywaniem incydentów z aplikacji                                    | string\[\]                          |
| sinceMs                  | Data określająca początek okresu analizowanego przez raport. Wyrażona w milisekundach od początku epoki | number                              |
| toMs                     | Data określająca koniec okresu analizowanego przez raport. Wyrażona w milisekundach od początku epoki   | number                              |
| clusterId                | Identyfikator klastra, z którego generowana jest **obserwacja**                                         | string                              |
| applicationConfiguration | Lista konfiguracji aplikacji, na podstawie których ma być wygenerowany raport                           | ApplicationInsightConfiguration\[\] |

#### ApplicationInsightConfiguration

Reprezentuje konfiguracje aplikacji wykorzystywaną podczas generowania raportu.

| Nazwa atrybutu  | Znaczenie                                                                               | Dziedzina |
| :-------------- | :-------------------------------------------------------------------------------------- | :-------- |
| applicationName | Nazwa aplikacji                                                                         | string    |
| accuracy        | Dokładność analizy                                                                      | string    |
| customPrompt    | Własne wejście do modelu językowego przy generowaniu **obserwacji** dla danej aplikacji | string    |

#### ApplicationIncidentSource

Reprezentuje źródło incydentu aplikacji. Zawiera metadane i zawartość loga, na podstawie którego został wykryty incydent.

| Nazwa atrybutu | Znaczenie                                         | Dziedzina |
| :------------- | :------------------------------------------------ | :-------- |
| timestamp      | Czas zebrania źródła incydentu                    | string    |
| podName        | Nazwa poda w której znajdowała się aplikacja      | string    |
| containerName  | Nazwa kontenera w którym znajdowała się aplikacja | string    |
| content        | Zawartość loga                                    | string    |
| image          | Identyfikator obrazu kontenera                    | string    |

#### NodeReport

Przechowuje incydenty występujące dla konkretnego hosta w ramach **Report**. Pozwala na połączenie informacji dotyczącej konfiguracji hosta przy generowaniu raportu wraz z incydentami z danego hosta.

| Nazwa atrybutu | Znaczenie                                                                                           | Dziedzina         |
| :------------- | :-------------------------------------------------------------------------------------------------- | :---------------- |
| node           | Nazwa hosta                                                                                         | string            |
| accuracy       | Dokładność raportu aplikacji                                                                        | string            |
| customPrompt   | Własne dodatkowe wejście do modelu językowego podczas interpretacji logów z hostów w ramach raportu | string            |
| incidents      | Lista incidentów z hostów                                                                           | NodeIncidents\[\] |

#### NodeIncident

Incydent hosta wraz z metadanymi pozwalającymi na stwierdzenie na podstawie których logów dany incydent został wykryty. Dokument ten został zdenormalizowany (posiada źródła incydentu (**NodeIncidentSources**) oraz nazwę hosta i konfigurację hosta, ponieważ pobieranie pojedyneczego incydentu wraz z jego metadanymi jest typowym zapytaniem, a denormalizacja pozwala na uniknięcie częstej operacji złączania dokumentów.

| Nazwa atrybutu | Znaczenie                                                                                              | Dziedzina              |
| :------------- | :----------------------------------------------------------------------------------------------------- | :--------------------- |
| id             | Unikalny identyfikator incydentu                                                                       | string                 |
| title          | Tytuł incydentu                                                                                        | string                 |
| customPrompt   | Własne dodatkowe wejście do modelu językowego podczas interpretacji logów z aplikacji w ramach raportu | string                 |
| clusterId      | Identyfikator klastra na którym wystąpił dany incydent                                                 | string                 |
| nodeName       | Nazwa hosta dla któregoj wystąpił dany incydent                                                        | string                 |
| category       | Kategoria incydentu                                                                                    | string                 |
| summary        | Podsumowanie incydentu                                                                                 | string                 |
| recommendation | Rekomendacja odnośnie rozwiązania incydentu                                                            | string                 |
| urgency        | Krytyczność incydentu                                                                                  | string                 |
| sources        | Źródła na podstawie, których dany incydent był wykryty                                                 | NodeIncidentSource\[\] |

#### ScheduledNodeInsights

Reprezentuje wykonywane zadanie generowania **obserwacji** z hostów.

| Nazwa atrybutu    | Znaczenie                                                                                               | Dziedzina                    |
| :---------------- | :------------------------------------------------------------------------------------------------------ | :--------------------------- |
| scheduledJobIds   | Identyfikatory zadań związanych z wykrywaniem incydentów z hosta                                        | string\[\]                   |
| sinceMs           | Data określająca początek okresu analizowanego przez raport. Wyrażona w milisekundach od początku epoki | number                       |
| toMs              | Data określająca koniec okresu analizowanego przez raport. Wyrażona w milisekundach od początku epoki   | number                       |
| clusterId         | Identyfikator klastra, z którego generowany jest raport                                                 | string                       |
| nodeConfiguration | Lista konfiguracji aplikacji, na podstawie których ma być wygenerowany raport                           | NodeInsightConfiguration\[\] |

#### NodeInsightConfiguration

Reprezentuje konfiguracje hosta wykorzystywaną podczas generowania raportu.

| Nazwa atrybutu | Znaczenie                                                                            | Dziedzina |
| :------------- | :----------------------------------------------------------------------------------- | :-------- |
| nodeName       | Nazwa hosta                                                                          | string    |
| accuracy       | Dokładność analizy                                                                   | string    |
| customPrompt   | Własne wejście do modelu językowego przy generowaniu **obserwacji** dla danego hosta | string    |

#### NodeIncidentSource

Reprezentuje źródło incydentu hosta. Zawiera metadane i zawartość loga, na podstawie którego został wykryty incydent.

| Nazwa atrybutu | Znaczenie                              | Dziedzina |
| :------------- | :------------------------------------- | :-------- |
| timestamp      | Czas zebrania źródła incydentu         | string    |
| filename       | Nazwa pliku, z którego zebrany był log | string    |
| content        | Zawartość loga                         | string    |

#### ScheduledIncidentMergerJob

Reprezentuje wykonywaną operacje scalania incydentów. Kolekcja ta jest abstrakcją pozwalająca na powiązanie tego zadania z kolekcjami związanymi bezpośrednio z dostawcą rozwiązania oferującego scalanie incydentów (np. model językowy od OpenAI i **ScheduledOpenAiJob**). Identyfikator dokumentu tej kolekcji odpowiada identyfikatorowi dokumentu posiadającego dane o bezpośredniej implementacji zadania scalania incydentów w innej kolekcji.

| Nazwa atrybutu | Znaczenie                                          | Dziedzina |
| :------------- | :------------------------------------------------- | :-------- |
| id             | Unikalny identyfikator zadania scalania incydentów | string    |

#### ScheduledOpenAiJob

Kolekcja odpowiadająca za przechowywanie zadań przekazywanych do modelu językowego od OpenAI.

| Nazwa atrybutu    | Znaczenie                                                                        | Dziedzina |
| :---------------- | :------------------------------------------------------------------------------- | :-------- |
| id                | Unikalny identyfikator                                                           | string    |
| scheduledAt       | Data zaplanowania wykonania zadania. Wyrażona w milisekundach od początku epoki. | number    |
| completionRequest | Zserializowane żądanie do Completion API od OpenAI                               | string    |
| status            | Status wykonania zadania                                                         | string    |
| batchId           | Zewnętrzny Identyfikator “batcha” od OpenAI                                      | string    |

<figure>
    <img src="/reports/reports-reportsdb.svg">
    <figcaption> Reports Service: Diagram bazy raportów [źródło opracowanie własne]</figcaption>
</figure>

### 7.3.3 Baza danych metadata service {#bazy-danych-metadata-service}

Metadata Service przechowuje otrzymane z brokera Kafki metadane hostów oraz klastrów, a także ich zagregowane wersje.

<figure>
    <img src="/metadata-service/database/metadata-service-database.svg">
    <figcaption> Management Service: Diagram bazy raportów [źródło opracowanie własne]</figcaption>
</figure>

**ApplicationMetadata**

Przechowuje metadane aplikacji.

| Nazwa atrybutu | Znaczenie                                                          | Dziedzina   |
| :------------- | :----------------------------------------------------------------- | :---------- |
| id             | Unikalny identyfikator dokumentu w ramach danego indeksu           | string      |
| clusterId      | Identyfikator klastra, na którym działa aplikacja                  | string      |
| collectedAtMs  | Czas kiedy metadane były zebrane w milisekundach od początku epoki | int64       |
| applications   | Aplikacje wchodzące w skład metadanych                             | Application |

**Application**

Przechowuje metadane konkretnej aplikacji.

| Nazwa atrybutu | Znaczenie                           | Dziedzina |
| :------------- | :---------------------------------- | :-------- |
| kind           | Identyfikator zasobu w Kubernetesie | string    |
| name           | Nazwa aplikacji                     | string    |

**NodeMetadata**

Przechowuje metadane hostów.

| Nazwa atrybutu | Znaczenie                                                          | Dziedzina  |
| :------------- | :----------------------------------------------------------------- | :--------- |
| id             | Unikalny identyfikator dokumentu w ramach danego indeksu           | string     |
| nodeName       | Nazwa hosta                                                        | string     |
| clusterId      | Identyfikator klastra, na którym działa aplikacja                  | string     |
| collectedAtMs  | Czas kiedy metadane były zebrane w milisekundach od początku epoki | int64      |
| watchedFiles   | Pliki na hoście, z których zbierane są logi                        | string\[\] |

**AggregatedApplicationMetadata**

Przechowuje zagregowane metadane aplikacji.

| Nazwa atrybutu | Znaczenie                                                          | Dziedzina |
| :------------- | :----------------------------------------------------------------- | :-------- |
| id             | Unikalny identyfikator dokumentu w ramach danego indeksu           | string    |
| clusterId      | Identyfikator klastra, na którym działa aplikacja                  | string    |
| collectedAtMs  | Czas kiedy metadane były zebrane w milisekundach od początku epoki | int64     |
| metadata       | Aplikacje wchodzące w skład metadanych                             | Metadata  |

**Metadata**

Przechowuje metadane aplikacji.

| Nazwa atrybutu | Znaczenie                   | Dziedzina |
| :------------- | :-------------------------- | :-------- |
| kind           | Nazwa zasobu w Kubernetesie | string    |
| name           | Nazwa aplikacji             | string    |

**AggregatedNodeMetadata**

Przechowuje zagregowane metadane hostów.

| Nazwa atrybutu | Znaczenie                                                          | Dziedzina |
| :------------- | :----------------------------------------------------------------- | :-------- |
| id             | Unikalny identyfikator dokumentu w ramach danego indeksu           | string    |
| clusterId      | Identyfikator klastra, na którym działa aplikacja                  | string    |
| collectedAtMs  | Czas kiedy metadane były zebrane w milisekundach od początku epoki | int64     |
| metadata       | Hosty wchodzące w skład metadanych                                 | Metadata  |

**Metadata**

Przechowuje metadane konkretnych hostów.

| Nazwa atrybutu | Znaczenie                                   | Dziedzina  |
| :------------- | :------------------------------------------ | :--------- |
| name           | Nazwa hosta                                 | string     |
| files          | Pliki na hoście, z których zbierane są logi | string\[\] |

**AggregatedClusterState**

Przechowuje zagregowane metadane klastrów.

| Nazwa atrybutu | Znaczenie                                                          | Dziedzina |
| :------------- | :----------------------------------------------------------------- | :-------- |
| id             | Unikalny identyfikator dokumentu w ramach danego indeksu           | string    |
| clusterId      | Identyfikator klastra, na którym działa aplikacja                  | string    |
| collectedAtMs  | Czas kiedy metadane były zebrane w milisekundach od początku epoki | int64     |
| metadata       | Klastry wchodzące w skład metadanych                               | Metadata  |

**Metadata**

Przechowuje metadane konkretnych klastrów.

| Nazwa atrybutu | Znaczenie             | Dziedzina |
| :------------- | :-------------------- | :-------- |
| clusterId      | Identyfikator klastra | string    |

### 7.3.4 Bazy danych management-service {#bazy-danych-management-service}

Management Service przechowuje swoją kopię wygenerowanych raportów oraz otrzymanych metadanych.

#### 7.3.4.1 Baza danych raportów

Raporty przechowywane są w schemacie lustrzanym do bazy danych Report Service, różnicą są kolekcje **ApplicationIncident**, **NodeIncident**, **ApplicationIncidentSource** oraz **NodeIncidentSource**, które denormalizują schemat raportu w celu szybszego dostępu do danych.

<figure>
    <img src="/management-service/database/management-service-database-mongodb-reports.svg">
    <figcaption> Management Service: Diagram bazy raportów [źródło opracowanie własne]</figcaption>
</figure>

**Report**

Obiekt ten przechowuje raport otrzymany od **Report Service**.

| Nazwa atrybutu                         | Znaczenie                                                                                               | Dziedzina                                 |
| :------------------------------------- | :------------------------------------------------------------------------------------------------------ | :---------------------------------------- |
| id                                     | Unikalny identyfikator raportu                                                                          | string                                    |
| correlationId                          | Identyfikator korelacji (z żądania na podstawie, którego został wygenerowany raport)                    | string                                    |
| status                                 | Aktualny status raportu.                                                                                | string                                    |
| sinceMs                                | Data określająca początek okresu analizowanego przez raport. Wyrażona w milisekundach od początku epoki | number                                    |
| toMs                                   | Data określająca koniec okresu analizowanego przez raport. Wyrażona w milisekundach od początku epoki   | number                                    |
| requestedAt                            | Data zażądania raportu. Wyrażona w milisekundach od początku epoki                                      | number                                    |
| scheduledGenerationAtMs                | Data wygenerowania raportu. Wyrażona w milisekundach od początku epoki                                  | number                                    |
| title                                  | Tytuł raportu.                                                                                          | string                                    |
| nodeReports                            | Incydenty z hostów pogrupowane po hostcie i konfiguracji z jakimi został wygenerowany.                  | NodeReport\[\]                            |
| applicationReports                     | Incydenty z aplikacji pogrupowane po aplikacji i konfiguracji z jakimi został wygenerowany.             | ApplicationReport\[\]                     |
| totalApplicationEntries                | Liczba logów z aplikacji, uwzględnionych podczas generowani raportu                                     | number                                    |
| totalNodeEntries                       | Liczba logów z hostów, uwzględnionych podczas generowani raportu                                        | number                                    |
| urgency                                | Poziom “krytyczności” całego raportu                                                                    | string                                    |
| scheduledApplicationInsights           | Zaplanowane do przetworzania (w ramach raportu) logi i konfiguracje aplikacji                           | ScheduledApplicationInsights\[\]          |
| scheduledNodeInsights                  | Zaplanowane do przetworzania (w ramach raportu) logi i konfiguracje hostów                              | ScheduledNodeInsights\[\]                 |
| analyzedApplications                   | Liczba aplikacji, przeanalizowanych w ramach raportu                                                    | number                                    |
| analyzedNodes                          | Liczba hostów, przeanalizowanych w ramach raportu                                                       | number                                    |
| scheduledApplicationIncidentMergerJobs | Zaplanowane zadania związane ze scalaniem zduplikowanych incydentów aplikacji                           | ScheduledApplicationIncidentMergerJob\[\] |
| scheduledNodeIncidentMergerJobs        | Zaplanowane zadania związane ze scalaniem zduplikowanych incydentów hostów                              | ScheduledNodeIncidentMergerJob\[\]        |

**ApplicationReport**

Przechowuje incydenty występujące dla konkretnej aplikacji w ramach **Report**. Pozwala na połączenie informacji dotyczącej konfiguracji aplikacji przy generowaniu raportu wraz z incydentami z danej aplikacji.

| Nazwa atrybutu  | Znaczenie                                                                                              | Dziedzina                |
| :-------------- | :----------------------------------------------------------------------------------------------------- | :----------------------- |
| applicationName | Nazwa aplikacji                                                                                        | string                   |
| accuracy        | Dokładność raportu aplikacji                                                                           | string                   |
| customPrompt    | Własne dodatkowe wejście do modelu językowego podczas interpretacji logów z aplikacji w ramach raportu | string                   |
| incidents       | Lista incidentów z aplikacji                                                                           | ApplicationIncidents\[\] |

**ApplicationIncident**

Incydent aplikacji wraz z metadanymi pozwalającymi na stwierdzenie na podstawie których logów dany incydent został wykryty. Dokument ten został zdenormalizowany (posiada źródła incydentu (**ApplicationIncidentSources**) oraz nazwę aplikacji i konfigurację aplikacji, ponieważ pobieranie pojedyneczego incydentu wraz z jego metadanymi jest typowym zapytaniem, a denormalizacja pozwala na uniknięcie częstej operacji złączania dokumentów.

| Nazwa atrybutu  | Znaczenie                                                                                              | Dziedzina                     |
| :-------------- | :----------------------------------------------------------------------------------------------------- | :---------------------------- |
| id              | Unikalny identyfikator incydentu                                                                       | string                        |
| title           | Tytuł incydentu                                                                                        | string                        |
| customPrompt    | Własne dodatkowe wejście do modelu językowego podczas interpretacji logów z aplikacji w ramach raportu | string                        |
| clusterId       | Identyfikator klastra na którym wystąpił dany incydent                                                 | string                        |
| applicationName | Nazwa aplikacji dla której wystąpił dany incydent                                                      | string                        |
| category        | Kategoria incydentu                                                                                    | string                        |
| summary         | Podsumowanie incydentu                                                                                 | string                        |
| recommendation  | Rekomendacja odnośnie rozwiązania incydentu                                                            | string                        |
| urgency         | Krytyczność incydentu                                                                                  | string                        |
| sources         | Źródła na podstawie, których dany incydent był wykryty                                                 | ApplicationIncidentSource\[\] |

**ScheduledApplicationInsights**

Reprezentuje wykonywane zadanie generowania **obserwacji** z aplikacji.

| Nazwa atrybutu           | Znaczenie                                                                                               | Dziedzina                           |
| :----------------------- | :------------------------------------------------------------------------------------------------------ | :---------------------------------- |
| scheduledJobIds          | Identyfikatory zadań związanych z wykrywaniem incydentów z aplikacji                                    | string\[\]                          |
| sinceMs                  | Data określająca początek okresu analizowanego przez raport. Wyrażona w milisekundach od początku epoki | number                              |
| toMs                     | Data określająca koniec okresu analizowanego przez raport. Wyrażona w milisekundach od początku epoki   | number                              |
| clusterId                | Identyfikator klastra, z którego generowana jest **obserwacja**                                         | string                              |
| applicationConfiguration | Lista konfiguracji aplikacji, na podstawie których ma być wygenerowany raport                           | ApplicationInsightConfiguration\[\] |

**ApplicationInsightConfiguration**

Reprezentuje konfiguracje aplikacji wykorzystywaną podczas generowania raportu.

| Nazwa atrybutu  | Znaczenie                                                                               | Dziedzina |
| :-------------- | :-------------------------------------------------------------------------------------- | :-------- |
| applicationName | Nazwa aplikacji                                                                         | string    |
| accuracy        | Dokładność analizy                                                                      | string    |
| customPrompt    | Własne wejście do modelu językowego przy generowaniu **obserwacji** dla danej aplikacji | string    |

**ApplicationIncidentSource**

Reprezentuje źródło incydentu aplikacji. Zawiera metadane i zawartość loga, na podstawie którego został wykryty incydent.

| Nazwa atrybutu | Znaczenie                                         | Dziedzina |
| :------------- | :------------------------------------------------ | :-------- |
| timestamp      | Czas zebrania źródła incydentu                    | string    |
| podName        | Nazwa poda w której znajdowała się aplikacja      | string    |
| containerName  | Nazwa kontenera w którym znajdowała się aplikacja | string    |
| content        | Zawartość loga                                    | string    |
| image          | Identyfikator obrazu kontenera                    | string    |

**NodeReport**

Przechowuje incydenty występujące dla konkretnego hosta w ramach **Report**. Pozwala na połączenie informacji dotyczącej konfiguracji hosta przy generowaniu raportu wraz z incydentami z danego hosta.

| Nazwa atrybutu | Znaczenie                                                                                              | Dziedzina         |
| :------------- | :----------------------------------------------------------------------------------------------------- | :---------------- |
| node           | Nazwa hosta                                                                                            | string            |
| accuracy       | Dokładność raportu aplikacji                                                                           | string            |
| customPrompt   | Własne dodatkowe wejście do modelu językowego podczas interpretacji logów z aplikacji w ramach raportu | string            |
| incidents      | Lista incidentów z hostów                                                                              | NodeIncidents\[\] |

**NodeIncident**

Incydent hosta wraz z metadanymi pozwalającymi na stwierdzenie na podstawie których logów dany incydent został wykryty. Dokument ten został zdenormalizowany (posiada źródła incydentu (**NodeIncidentSources**) oraz nazwę hosta i konfigurację hosta, ponieważ pobieranie pojedyneczego incydentu wraz z jego metadanymi jest typowym zapytaniem, a denormalizacja pozwala na uniknięcie częstej operacji złączania dokumentów.

| Nazwa atrybutu | Znaczenie                                                                                              | Dziedzina              |
| :------------- | :----------------------------------------------------------------------------------------------------- | :--------------------- |
| id             | Unikalny identyfikator incydentu                                                                       | string                 |
| title          | Tytuł incydentu                                                                                        | string                 |
| customPrompt   | Własne dodatkowe wejście do modelu językowego podczas interpretacji logów z aplikacji w ramach raportu | string                 |
| clusterId      | Identyfikator klastra na którym wystąpił dany incydent                                                 | string                 |
| nodeName       | Nazwa hosta dla któregoj wystąpił dany incydent                                                        | string                 |
| category       | Kategoria incydentu                                                                                    | string                 |
| summary        | Podsumowanie incydentu                                                                                 | string                 |
| recommendation | Rekomendacja odnośnie rozwiązania incydentu                                                            | string                 |
| urgency        | Krytyczność incydentu                                                                                  | string                 |
| sources        | Źródła na podstawie, których dany incydent był wykryty                                                 | NodeIncidentSource\[\] |

**ScheduledNodeInsights**

Reprezentuje wykonywane zadanie generowania **obserwacji** z hostów.

| Nazwa atrybutu    | Znaczenie                                                                                               | Dziedzina                    |
| :---------------- | :------------------------------------------------------------------------------------------------------ | :--------------------------- |
| scheduledJobIds   | Identyfikatory zadań związanych z wykrywaniem incydentów z hosta                                        | string\[\]                   |
| sinceMs           | Data określająca początek okresu analizowanego przez raport. Wyrażona w milisekundach od początku epoki | number                       |
| toMs              | Data określająca koniec okresu analizowanego przez raport. Wyrażona w milisekundach od początku epoki   | number                       |
| clusterId         | Identyfikator klastra, z którego generowany jest raport                                                 | string                       |
| nodeConfiguration | Lista konfiguracji aplikacji, na podstawie których ma być wygenerowany raport                           | NodeInsightConfiguration\[\] |

**NodeInsightConfiguration**

Reprezentuje konfiguracje hosta wykorzystywaną podczas generowania raportu.

| Nazwa atrybutu | Znaczenie                                                                            | Dziedzina |
| :------------- | :----------------------------------------------------------------------------------- | :-------- |
| nodeName       | Nazwa aplikacji                                                                      | string    |
| accuracy       | Dokładność analizy                                                                   | string    |
| customPrompt   | Własne wejście do modelu językowego przy generowaniu **obserwacji** dla danego hosta | string    |

**NodeIncidentSource**

Reprezentuje źródło incydentu hosta. Zawiera metadane i zawartość loga, na podstawie którego został wykryty incydent.

| Nazwa atrybutu | Znaczenie                              | Dziedzina |
| :------------- | :------------------------------------- | :-------- |
| timestamp      | Czas zebrania źródła incydentu         | string    |
| filename       | Nazwa pliku, z którego zebrany był log | string    |
| content        | Zawartość loga                         | string    |
| image          | Identyfikator obrazu kontenera         | string    |

#### 7.3.4.2 Baza danych metadanych

Metadane o aplikacjach, hostach oraz klastrach są przechowywane w sposób lustrzany do Metadata Service. Rozszerzeniem jest natomiast zdenormaliowana kolekcja **ClusterHistory**, zawierająca informacje o wszystkich działających w obrębie danego klastra aplikacjach oraz hostach. Dodatkowo, występuje kolekcja **ReportGenerationRequestMetadata**, przechowująca dane audytowe o zapytaniach generacji raportu oraz ich statusie.

<figure>
    <img src="/management-service/database/management-service-database-mongodb-metadata-configuration-transparent.svg">
    <figcaption> Management Service: Diagram bazy metadanych [źródło opracowanie własne]</figcaption>
</figure>

**AggregatedApplicationMetadata**

Przechowuje zagregowane metadane aplikacji.

| Nazwa atrybutu | Znaczenie                                                          | Dziedzina |
| :------------- | :----------------------------------------------------------------- | :-------- |
| id             | Unikalny identyfikator dokumentu w ramach danego indeksu           | string    |
| clusterId      | Identyfikator klastra, na którym działa aplikacja                  | string    |
| collectedAtMs  | Czas kiedy metadane były zebrane w milisekundach od początku epoki | int64     |
| metadata       | Aplikacje wchodzące w skład metadanych                             | Metadata  |

**Metadata**

Przechowuje metadane o konkretnej aplikacji.

| Nazwa atrybutu | Znaczenie                   | Dziedzina |
| :------------- | :-------------------------- | :-------- |
| kind           | Nazwa zasobu w Kubernetesie | string    |
| name           | Nazwa aplikacji             | string    |

**AggregatedNodeMetadata**

Przechowuje zagregowane metadane hostów.

| Nazwa atrybutu | Znaczenie                                                          | Dziedzina |
| :------------- | :----------------------------------------------------------------- | :-------- |
| id             | Unikalny identyfikator dokumentu w ramach danego indeksu           | string    |
| clusterId      | Identyfikator klastra, na którym działa aplikacja                  | string    |
| collectedAtMs  | Czas kiedy metadane były zebrane w milisekundach od początku epoki | int64     |
| metadata       | Hosty wchodzące w skład metadanych                                 | Metadata  |

**Metadata**

Przechowuje metadane o konkretnym hoście.

| Nazwa atrybutu | Znaczenie                                   | Dziedzina  |
| :------------- | :------------------------------------------ | :--------- |
| name           | Nazwa hosta                                 | string     |
| files          | Pliki na hoście, z których zbierane są logi | string\[\] |

**AggregatedClusterState**

Przechowuje zagregowane metadane klastrów.

| Nazwa atrybutu | Znaczenie                                                          | Dziedzina |
| :------------- | :----------------------------------------------------------------- | :-------- |
| id             | Unikalny identyfikator dokumentu w ramach danego indeksu           | string    |
| clusterId      | Identyfikator klastra, na którym działa aplikacja                  | string    |
| collectedAtMs  | Czas kiedy metadane były zebrane w milisekundach od początku epoki | int64     |
| metadata       | Klastry wchodzące w skład metadanych                               | Metadata  |

**Metadata**

Przechowuje metadane o konkretnym klastrze.

| Nazwa atrybutu | Znaczenie             | Dziedzina |
| :------------- | :-------------------- | :-------- |
| clusterId      | Identyfikator klastra | string    |

**ReportGenerationRequestMetadata**

Agreguje metadane o zapytaniu generującym raport.

| Nazwa atrybutu | Znaczenie                                     | Dziedzina           |
| :------------- | :-------------------------------------------- | :------------------ |
| id             | Unikalny identyfikator                        | string              |
| status         | Aktualny status generowanego raportu          | string              |
| reportType     | Typ raportu (cykliczny lub na żądanie)        | string              |
| request        | Zapytanie z konfiguracją generowanego raportu | CreateReportRequest |

**CreateReportRequest**

Przechowuje metadane o ciele zapytania generującego raport.

| Nazwa atrybutu            | Znaczenie                                                                                                 | Dziedzina                    |
| :------------------------ | :-------------------------------------------------------------------------------------------------------- | :--------------------------- |
| clusterId                 | Unikalny identyfikator klastra dla którego generowany jest raport                                         | string                       |
| accuracy                  | Dokładność generowanego raportu                                                                           | string                       |
| sinceMs                   | Początek przedziału logów branych pod uwagę podczas generowania raportu w milisekundach od początku epoki | int64                        |
| toMs                      | Koniec przedziału logów branych pod uwagę podczas generowania raportu w milisekundach od początku epoki   | int64                        |
| slackReceiverIds          | Identyfikatory odbiorników powiadomień Slack                                                              | \[\]int                      |
| discordReceiverIds        | Identyfikatory odbiorników powiadomień Discord                                                            | \[\]int                      |
| emailReceiverIds          | Identyfikatory odbiorników powiadomień Email                                                              | \[\]int                      |
| nodeConfigurations        | Konfiguracje hostów dla generowanego raportu                                                              | \[\]NodeConfiguration        |
| applicationConfigurations | Konfiguracje aplikacji dla generowanego raportu                                                           | \[\]ApplicationConfiguration |

**ApplicationConfiguration**

Przechowuje dane o konfiguracji danej aplikacji, która ma wystąpić w raporcie.

| Nazwa atrybutu  | Znaczenie                                                                                              | Dziedzina |
| :-------------- | :----------------------------------------------------------------------------------------------------- | :-------- |
| applicationName | Unikalny identyfikator                                                                                 | string    |
| customPrompt    | Własne dodatkowe wejście do modelu językowego podczas interpretacji logów z aplikacji w ramach raportu | string    |
| accuracy        | Typ raportu (cykliczny lub na żądanie)                                                                 | string    |

**NodeConfiguration**

Przechowuje dane o konfiguracji danego hosta, który ma wystąpić w raporcie.

| Nazwa atrybutu | Znaczenie                                                                                              | Dziedzina |
| :------------- | :----------------------------------------------------------------------------------------------------- | :-------- |
| nodeName       | Unikalny identyfikator                                                                                 | string    |
| customPrompt   | Własne dodatkowe wejście do modelu językowego podczas interpretacji logów z aplikacji w ramach raportu | string    |
| accuracy       | Typ raportu (cykliczny lub na żądanie)                                                                 | string    |

**ClusterHistory**

Przechowuje dane o aplikacjach oraz hostach należących aktualnie lub historycznie do klastra.

| Nazwa atrybutu | Znaczenie                                         | Dziedzina   |
| :------------- | :------------------------------------------------ | :---------- |
| id             | Unikalny identyfikator klastra                    | string      |
| applications   | Aplikacje które działają lub działały na klastrze | Application |
| nodes          | Hosty które należą lub należały do klastra        | Node        |

**Application**

Przechowuje dane dotyczące aplikacji która występuje lub występowała w klastrze.

| Nazwa atrybutu | Znaczenie                                           | Dziedzina |
| :------------- | :-------------------------------------------------- | :-------- |
| name           | Nazwa aplikacji                                     | string    |
| kind           | Rodzaj zasobu w Kubernetes                          | string    |
| running        | Wskazuje czy aplikacja aktualnie działa na klastrze | boolean   |

**Node**

Przechowuje dane dotyczące hosta, który występuje lub występował w klastrze.

| Nazwa atrybutu | Znaczenie                                        | Dziedzina |
| :------------- | :----------------------------------------------- | :-------- |
| name           | Nazwa hosta                                      | string    |
| running        | Wskazuje czy host jest aktualnie częścią klastra | boolean   |

#### 7.3.4.3 Baza danych użytkowników i konfiguracji {#baza-danych-użytkowników-i-konfiguracji}

#### 7.3.4.3.1 Model fizyczny bazy danych

<figure>
    <img src="/management-service/database/management-service-postgres-diagram.drawio.svg">
    <figcaption> Management Service: Diagram fizyczny bazy danych [źródło opracowanie własne]</figcaption>
</figure>

#### 7.3.4.3.1 Definicja schematów relacji

##### Tabela 1: \_user

**User**(<u>id</u>, email, nickname, password, provider)

| Nazwa atrybutu | Znaczenie                        | Dziedzina    | Unikalność | OBL(+) OPC(-) |
| -------------- | -------------------------------- | ------------ | ---------- | ------------- |
| id             | Identyfikator użytkownika        | bigint       | +          | +             |
| email          | Adres e-mail użytkownika         | varchar(255) | +          | +             |
| nickname       | Pseudonim użytkownika            | varchar(255) | -          | +             |
| password       | Hasło użytkownika                | varchar(255) | -          | +             |
| provider       | Dostawca usługi uwierzytelniania | varchar(255) | -          | +             |

**Klucze kandydujące**: id, email, nickname
**Klucz główny**: id
**Zależności funkcyjne**
&nbsp;&nbsp;&nbsp;id → email, nickname, password, provider
&nbsp;&nbsp;&nbsp;email → id, nickname, password, provider

##### Tabela 2: discord_receiver

**DiscordReceiver**(<u>id</u>, created_at, receiver_name, updated_at, webhook_url)

| Nazwa atrybutu | Znaczenie                       | Dziedzina    | Unikalność | OBL(+) OPC(-) |
| -------------- | ------------------------------- | ------------ | ---------- | ------------- |
| id             | Identyfikator odbiorcy Discorda | bigint       | +          | +             |
| created_at     | Data utworzenia                 | bigint       | -          | +             |
| receiver_name  | Nazwa odbiorcy                  | varchar(255) | -          | +             |
| updated_at     | Data ostatniej aktualizacji     | bigint       | -          | -             |
| webhook_url    | URL webhooka odbiorcy Discorda  | varchar(255) | -          | +             |

**Klucze kandydujące**: id, receiver_name
**Klucz główny**: id
**Zależności funkcyjne**
&nbsp;&nbsp;&nbsp;id → created_at, receiver_name, updated_at, webhook_url

##### Tabela 3: email_receiver

**EmailReceiver**(<u>id</u>, created_at, receiver_email, receiver_name, updated_at)

| Nazwa atrybutu | Znaczenie                     | Dziedzina    | Unikalność | OBL(+) OPC(-) |
| -------------- | ----------------------------- | ------------ | ---------- | ------------- |
| id             | Identyfikator odbiorcy e-mail | bigint       | +          | +             |
| created_at     | Data utworzenia               | bigint       | -          | +             |
| receiver_email | Adres e-mail odbiorcy         | varchar(255) | -          | +             |
| receiver_name  | Nazwa odbiorcy                | varchar(255) | -          | +             |
| updated_at     | Data ostatniej aktualizacji   | bigint       | -          | -             |

**Klucze kandydujące**: id, receiver_email, receiver_name
**Klucz główny**: id
**Zależności funkcyjne**
&nbsp;&nbsp;&nbsp;id → created_at, receiver_email, updated_at, receiver_name

##### Tabela 4: slack_receiver

**SlackReceiver**(<u>id</u>, created_at, receiver_name, updated_at, webhook_url)

| Nazwa atrybutu | Znaczenie                     | Dziedzina    | Unikalność | OBL(+) OPC(-) |
| -------------- | ----------------------------- | ------------ | ---------- | ------------- |
| id             | Identyfikator odbiorcy Slacka | bigint       | +          | +             |
| created_at     | Data utworzenia               | bigint       | -          | +             |
| receiver_name  | Nazwa odbiorcy                | varchar(255) | -          | +             |
| updated_at     | Data ostatniej aktualizacji   | bigint       | -          | -             |
| webhook_url    | URL webhooka odbiorcy Slacka  | varchar(255) | -          | +             |

**Klucze kandydujące**: id, receiver_name
**Klucz główny**: id
**Zależności funkcyjne**
&nbsp;&nbsp;&nbsp;id → created_at, receiver_name, updated_at, webhook_url

##### Tabela 5: application_configuration

**ApplicationConfiguration**(<u>id</u>, accuracy, custom_prompt, kind, name)

| Nazwa atrybutu | Znaczenie               | Dziedzina    | Unikalność | OBL(+) OPC(-) |
| -------------- | ----------------------- | ------------ | ---------- | ------------- |
| id             | Identyfikator aplikacji | bigint       | +          | +             |
| accuracy       | Dokładność              | smallint     | -          | +             |
| custom_prompt  | Własny prompt           | varchar(255) | -          | -             |
| kind           | Typ aplikacji           | varchar(255) | -          | +             |
| name           | Nazwa aplikacji         | varchar(255) | -          | +             |

**Klucze kandydujące**: id, name
**Klucz główny**: id
**Zależności funkcyjne**
&nbsp;&nbsp;&nbsp;id → accuracy, custom_prompt, kind, name

##### Tabela 6: node_configuration

**NodeConfiguration**(<u>id</u>, name, accuracy, custom_prompt)

| Nazwa atrybutu | Znaczenie                       | Dziedzina    | Unikalność | OBL(+) OPC(-) |
| -------------- | ------------------------------- | ------------ | ---------- | ------------- |
| id             | Identyfikator konfiguracji noda | bigint       | +          | +             |
| name           | Nazwa noda                      | varchar(255) | +          | +             |
| accuracy       | Dokładność                      | bigint       | -          | +             |
| custom_prompt  | Własny prompt                   | bigint       | -          | -             |

**Klucze kandydujące**: id, name
**Klucz główny**: id
**Zależności funkcyjne**
&nbsp;&nbsp;&nbsp;id → name, accuracy, custom_prompt

##### Tabela 7: cluster_configuration

ClusterConfiguration(<u>id</u>, accuracy, generated_every_millis, is_enabled)

| Nazwa atrybutu         | Znaczenie                          | Dziedzina    | Unikalność | OBL(+) OPC(-) |
| ---------------------- | ---------------------------------- | ------------ | ---------- | ------------- |
| id                     | Identyfikator konfiguracji klastra | varchar(255) | +          | +             |
| accuracy               | Dokładność                         | smallint     | -          | +             |
| generated_every_millis | Okres generowania w milisekundach  | bigint       | -          | +             |
| is_enabled             | Status konfiguracji                | boolean      | -          | +             |

**Klucze kandydujące**: id
**Klucz główny**: id
**Zależności funkcyjne**
&nbsp;&nbsp;&nbsp;id → accuracy, generated_every_millis, is_enabled

##### Tabela 8: cluster_schedule

**ClusterSchedule**(<u>cluster_id</u>, last_generation_ms, period_ms)

| Nazwa atrybutu     | Znaczenie                                | Dziedzina    | Unikalność | OBL(+) OPC(-) |
| ------------------ | ---------------------------------------- | ------------ | ---------- | ------------- |
| cluster_id         | Identyfikator klastra                    | varchar(255) | +          | +             |
| last_generation_ms | Ostatnia czas wygenerowania raportu w ms | bigint       | -          | +             |
| period_ms          | Przedział czasu miedzy raportami         | bigint       | -          | +             |

**Klucze kandydujące**: cluster_id
**Klucz główny**: cluster_id
**Zależności funkcyjne**
&nbsp;&nbsp;&nbsp;cluster_id → last_generation_ms, period_ms

## 7.4. Interfejsy programistyczne {#interfejsy-programistyczne}

### Żądanie raportu - ReportRequested (Reports Service):

Schemat wiadomości z żądaniem wygenerowania raportu zawiera informacje o tym z jakiego klastra i przedziału czasu powinien być raport. Dodatkowo żądanie takie posiada konfiguracje aplikacji i hostów, z którą raport powinien być wygenerowany. Każda taka konfiguracja posiada nazwę aplikacji (jako jej identyfikator) oraz dokładność analizy logów z danej aplikacji oraz własne dodatkowe wejście do modelu językowego wykorzystywanego do generowania raprotu.

Ponadto, aby nadawca, mógł powiązać żądanie z odpowiedzią jaką dostanie, wysyłany jest identyfikator korelacji.

<figure>
    <img src="/asyncapi-screens/report-requested.png">
    <figcaption> ReportRequested: Zrzut ekranu z AsyncAPI [źródło opracowanie własne]</figcaption>
</figure>

Kanał odpowiedzialny za przesyłanie konfiguracji raportów do wygenerowania.

### Wygenerowany raport - ReportGenerated (Reports Service):

Wiadomość z wygenerowanym raportem, zawiera nie tylko samą encje raportu z incydentami oraz metadanymi związanymi z czasem i konfiguracją raportu, ale również identyfikator korelacji, który odpowiada temu, z którym dany raport był zażądany.

<figure>
    <img src="/asyncapi-screens/report-generated.png">
    <figcaption> ReportGenerated: Zrzut ekranu z AsyncAPI [źródło opracowanie własne]</figcaption>
</figure>

Kanał odpowiedzialny za przesyłanie gotowych raportów

### Błąd podczas przetwarzania żądania wygenerowania raportu - ReportRequestFailed (Reports Service):

Wiadomość ta jest wysyłana podczas jakiegokolwiek krytycznego błędów podczas generowania raportu. Wiadomość taka posiada identyfikator korelacji odpowiadający żądaniu wygenerowania raportu, z którego przetwarzaniem wystąpił problem. Dodatkowo wiadomość taka posiada typ, czas wystąpienia oraz treść błędu.

<figure>
    <img src="/asyncapi-screens/report-request-failed.png">
    <figcaption> ReportRequestFailed: Zrzut ekranu z AsyncAPI [źródło opracowanie własne]</figcaption>
</figure>

### Logi z hostów - NodeLogs (Logs Ingestion Service)

Wiadomość zawierająca logi z hostów wraz z metadanymi o hoście, dacie wyprodukowania danego loga i jego zawartość.

<figure>
    <img src="/logs-ingestion/logs-ingestion-node-logs-asyncapi.png">
    <figcaption> NodeLogs: Zrzut ekranu z AsyncAPI [źródło opracowanie własne]</figcaption>
</figure>

### Logi z aplikacji - ApplicationLogs (Logs Ingestion Service)

Wiadomość zawierająca logi z aplikacji wraz z metadanymi o aplikacji, dacie wyprodukowania danego loga i jego zawartość.

<figure>
    <img src="/logs-ingestion/logs-ingestion-applicatino-logs-asyncapi.png">
    <figcaption> ApplicationLogs: Zrzut ekranu z AsyncAPI [źródło opracowanie własne]</figcaption>
</figure>

### Przesyłanie zagregowanych metadanych aplikacji - ApplicationMetadataUpdated (Metadata Service):

Zmiana stanu zagregowanych metadanych aplikacji jest emitowana przez Metadata Service.

<figure>
    <img src="/asyncapi-screens/application-metadata-updated-pub.png">
    <figcaption> ApplicationMetadataUpdated: Zrzut ekranu z AsyncAPI [źródło opracowanie własne]</figcaption>
</figure>

<figure>
    <img src="/asyncapi-screens/application-metadata-updated-sub.png">
    <figcaption> ApplicationMetadataUpdated: Zrzut ekranu z AsyncAPI  [źródło opracowanie własne]</figcaption>
</figure>

### Przesyłanie zagregowanych metadanych hostów - NodeMetadataUpdated (Metadata Service):

Zmiana stanu zagregowanych metadanych hostów jest emitowana przez Metadata Service.

<figure>
    <img src="/asyncapi-screens/node-metadata-updated-pub.png">
    <figcaption> NodeMetadataUpdated: Zrzut ekranu z AsyncAPI [źródło opracowanie własne]</figcaption>
</figure>

<figure>
    <img src="/asyncapi-screens/node-metadata-updated-sub.png">
    <figcaption> NodeMetadataUpdated: Zrzut ekranu z AsyncAPI [źródło opracowanie własne]</figcaption>
</figure>

### Przesyłanie zagregowanych metadanych klastrów - ClusterMetadataUpdated (Metadata Service):

Zmiana stanu zagregowanych metadanych klastrów jest emitowana przez Metadata Service.

<figure>
    <img src="/asyncapi-screens/cluster-metadata-updated-pub.png">
    <figcaption> ClusterMetadataUpdated: Zrzut ekranu z AsyncAPI [źródło opracowanie własne]</figcaption>
</figure>

<figure>
    <img src="/asyncapi-screens/cluster-metadata-updated-sub.png">
    <figcaption> ClusterMetadataUpdated: Zrzut ekranu z AsyncAPI [źródło opracowanie własne]</figcaption>
</figure>

## 7.5 Projekt interfejsu {#projekt-interfejsu}

Na podstawie wymagań funkcjonalnych oraz historyjek użytkownika zaprojektowano interfejs, który odpowiada oczekiwaniom użytkowników podczas korzystania z systemu. Szczególną uwagę poświęcono zapewnieniu pozytywnych doświadczeń użytkownika z Magpie Monitorem. W tym celu skupiono się na prostocie, czytelności oraz estetycznym wyglądzie interfejsu. Dodatkowym wyzwaniem było zapewnienie pełnej responsywności, aby użytkownicy mogli wygodnie korzystać z systemu także na urządzeniach mobilnych za pośrednictwem przeglądarki.

### 7.5.1 Widok logowania {#widok-logowania}

<figure>
    <img src="/user-interface/login-page.png">
    <figcaption> Widok logowania [źródło opracowanie własne]</figcaption>
</figure>

<figure>
    <img src="/user-interface/login-page-mobile.png">
    <figcaption> Widok logowania na urządzeniu mobilnym [źródło opracowanie własne]</figcaption>
</figure>

Pierwszym widokiem wyświetlanym po wejściu na witrynę Magpie Monitor jest ekran logowania. Jego celem jest spełnienie wymogu uwierzytelnienia użytkownika przed uzyskaniem dostępu do funkcji systemu. Po wybraniu opcji „Sign in with Google” użytkownik zostaje przekierowany na ekran logowania dostarczany przez firmę Google, gdzie może wybrać konto, które chce wykorzystać do zalogowania się do systemu.

### 7.5.2 Widok główny {#widok-główny}

<figure>
    <img src="/user-interface/main-page.png">
    <figcaption> Widok główny [źródło opracowanie własne]</figcaption>
</figure>

Po zalogowaniu się do systemu użytkownik zostaje przeniesiony na stronę główną. Wyświetla ona ostatnio wygenerowany raport, wraz z informacją o klastrze Kubernetes, z którego pochodzi, oraz przedziałem czasowym, w którym zebrano logi wykorzystane do jego stworzenia.

W sekcji **“Statistics”** prezentowane są użytkownikowi kluczowe dane opisane w historyjkach użytkownika, takie jak:

- liczba przeanalizowanych aplikacji,
- liczba przeanalizowanych hostów,
- liczba incydentów z podziałem według pilności, dodatkowo wizualizowana na wykresie kołowym,
- liczba przeanalizowanych logów pochodzących z aplikacji,
- liczba przeanalizowanych logów pochodzących z hostów,
- nazwa hosta z największą liczbą wykrytych incydentów oraz ich liczba,
- nazwa aplikacji z największą liczbą wykrytych incydentów oraz ich liczba.

Po lewej stronie znajduje się pionowy pasek nawigacyjny. W jego górnej części umieszczono logo systemu oraz jego nazwę. Kliknięcie logo z dowolnej podstrony przenosi użytkownika na stronę główną. Poniżej znajdują się sekcje ułatwiające nawigację po systemie. Obecnie aktywny widok jest wyróżniony kolorem zielonym, co zgodnie z pierwszą heurystyką Nielsena zapewnia użytkownikowi jasność dotyczącą aktualnego stanu systemu.

Na samym dole paska nawigacyjnego znajduje się opcja **wylogowania**, która przenosi użytkownika z powrotem na ekran logowania oraz odbiera mu uprawnienia do korzystania z systemu.

Po przewinięciu strony głównej w dół użytkownik zobaczy kolejny widok:

<figure>
    <img src="/user-interface/main-page-2.png">
    <figcaption> Widok główny kontynuacja [źródło opracowanie własne]</figcaption>
</figure>

Sekcje **“Application Incidents”** oraz **“Node Incidents”** zawierają listy incydentów wykrytych odpowiednio w aplikacjach i hostach. Każdy wiersz na liście przedstawia następujące informacje:

- nazwę hosta lub aplikacji, z której pochodzi incydent,
- kategorię incydentu,
- tytuł incydentu,
- datę wystąpienia incydentu.

Kliknięcie dowolnego wiersza z obu list przenosi użytkownika na stronę szczegółową, dedykowaną wybranemu incydentowi.

W wersji zoptymalizowanej dla urządzeń mobilnych strona jest odpowiednio dostosowana do mniejszych ekranów i prezentuje się w sposób przyjazny dla użytkownika:

<figure>
    <img src="/user-interface/main-page-mobile.png">
    <figcaption> Widok główny na urządzeniach mobilnych [źródło opracowanie własne]</figcaption>
</figure>

Pasek nawigacji został przeniesiony na górę strony, a zakładki są prezentowane wyłącznie za pomocą ikon. Ta zmiana ma na celu ułatwienie nawigacji w wersji mobilnej, zapewniając większą czytelność i oszczędność miejsca na mniejszych ekranach.
Oprócz tej modyfikacji układ strony pozostał w dużej mierze niezmieniony. Główną różnicą jest zwiększenie liczby wierszy, w których prezentowane są statystyki, co pozwala na bardziej efektywne wykorzystanie przestrzeni ekranu. Pozostałe elementy interfejsu nie zostały zmienione.

### 7.5.3 Widok incydentu {#widok-incydentu}

<figure>
    <img src="/user-interface/incident-page.png">
    <figcaption> Widok incydentu [źródło opracowanie własne]</figcaption>
</figure>

Widok incydentu ma na celu dostarczenie użytkownikowi szczegółowych informacji na temat wykrytego problemu. Górną część strony zajmuje nazwa incydentu, którą nadał model językowy, służąca jako zwięzłe określenie charakteru problemu. Obok znajduje się data, która wskazuje moment wystąpienia incydentu.

Bezpośrednio pod nazwą umieszczona jest sekcja **"Application"**, prezentująca informacje o nazwie aplikacji oraz klastrze, w którym wykryto problem. Sekcja ta zawiera również przedział czasowy, w którym system zebrał logi sugerujące opisywany incydent.

Po prawej stronie widoczna jest sekcja **"Configuration"**, przedstawiająca szczegóły konfiguracji raportu wybrane przez użytkownika.

W dolnej części strony znajdują się kolejne sekcje:

- **"Summary"** – zwięzły opis incydentu, ułatwiający użytkownikowi szybkie zrozumienie istoty problemu.
- **"Recommendation"** – zestaw sugerowanych działań mających na celu rozwiązanie wykrytego problemu.

Na samym dole użytkownik ma dostęp do logów, które wskazywały na wystąpienie incydentu, co umożliwia szczegółową analizę danych źródłowych.

<figure>
    <img src="/user-interface/incident-page-mobile.png">
    <figcaption> Widok incydentu na urządzeniach mobilnych [źródło opracowanie własne]</figcaption>
</figure>

W wersji mobilnej widok został dostosowany w sposób analogiczny do wersji desktopowej. Dla wygody użytkownika sekcje zostały ułożone w układ kolumnowy, co pozwala na przewijanie ekranu wyłącznie w pionie. Dzięki temu nawigacja staje się bardziej intuicyjna i dostosowana do korzystania na urządzeniach z mniejszym ekranem.

### 7.5.4 Widok raportów {#widok-raportów}

<figure>
    <img src="/user-interface/reports-page.png">
    <figcaption> Widok raportów [źródło opracowanie własne]</figcaption>
</figure>

Widok raportów składa się z listy raportów posortowanej malejąco według daty. Raporty są podzielone na dwie sekcje: **“Scheduled Reports”** i **“Reports on Demand”**. Każdy wiersz w liście zawiera następujące informacje:

- nazwę klastra, z którego pochodzi raport,
- tytuł raportu,
- poziom pilności,
- zakres czasu, w którym zebrano logi wykorzystane w raporcie,
- datę zatwierdzenia generacji raportu.

Kliknięcie nazwy klastra lub przycisku w kolumnie **“Actions”** przenosi użytkownika na szczegółowy widok raportu, który jest analogiczny do widoku głównego.

<figure>
    <img src="/user-interface/reports-page-mobile.png">
    <figcaption> Widok raportów na urządzeniach mobilnych [źródło opracowanie własne]</figcaption>
</figure>

Widok w wersji mobilnej nie różni się znacząco od wersji na komputery stacjonarne. Ze względu na ograniczoną szerokość ekranu, użytkownik musi przesuwać wiersz w lewo, aby zobaczyć pozostałe kolumny. Choć sprawia to, że korzystanie z widoku mobilnego jest mniej wygodne niż w przypadku wersji desktopowej, interfejs nadal pozostaje funkcjonalny i intuicyjny.

### 7.5.5 Widok klastrów {#widok-klastrów}

<figure>
    <img src="/user-interface/clusters-page.png">
    <figcaption> Widok klastrów [źródło opracowanie własne]</figcaption>
</figure>

Widok klastrów przedstawia listę klastrów, które mogą być monitorowane przez Magpie Monitor. Tabela zawiera kolumny opisujące konfigurację raportów dla danego wiersza. Kliknięcie nazwy lub przycisku „Reports Configuration” przenosi użytkownika do widoku konfiguracji raportów.

Widok ten na urządzeniach mobilnych jest analogiczny do widoku raportów. Ze względu na ograniczenia ekranu, użytkownik może przesuwać tabelę w poziomie, aby zobaczyć wszystkie kolumny, co zapewnia pełną funkcjonalność również na urządzeniach mobilnych.

### 7.5.6 Widok konfiguracji raportów {#widok-konfiguracji-raportów}

<figure>
    <img src="/user-interface/report-config-on-demand.png">
    <figcaption> Widok konfiguracji raportów [źródło opracowanie własne]</figcaption>
</figure>

Widok konfiguracji raportu zawiera następujące sekcje:

- “Default Accuracy” – określa domyślną dokładność, jaką będą miały nowo dodane aplikacje i hosty w raporcie.
- “Generation Type” – umożliwia wybór, czy raport ma być generowany na żądanie (dla wybranego przedziału czasu), czy cyklicznie (w określonych odstępach czasu, np. co tydzień lub co miesiąc).
- “Date Range” – w przypadku raportu na żądanie pozwala określić przedział czasu, z którego logi będą analizowane. Dla raportów cyklicznych umożliwia wybór częstotliwości ich generowania.
- “Notification Channels” – definiuje kanały powiadomień, na które mają być wysyłane informacje dotyczące raportu. Po kliknięciu ikony „plus” wyświetla się okienko z listą skonfigurowanych kanałów, które można dodać.
- “Application” – zawiera listę aplikacji uwzględnionych w raporcie. Użytkownik może zmieniać dokładność dla każdej z nich oraz dodawać indywidualne instrukcje wejściowe. Kliknięcie ikony „plus” otwiera okienko z listą dostępnych aplikacji do wyboru.
- “Nodes” – zawiera listę hostów uwzględnionych w raporcie. Funkcjonalność tej sekcji jest identyczna jak w przypadku sekcji "Application”.

W prawym dolnym rogu znajdują się dwa przyciski: jeden odpowiada za wygenerowanie raportu, a drugi za anulowanie zmian, co powoduje powrót do widoku głównego.  
W wersji mobilnej sekcje są ułożone w jednej kolumnie, co poprawia czytelność i ułatwia nawigację po stronie.

### 7.5.7 Widok konfiguracji powiadomień {#widok-konfiguracji-powiadomień}

<figure>
    <img src="/user-interface/notification-page.png">
    <figcaption> Widok kanałów powiadomień [źródło opracowanie własne]</figcaption>
</figure>

Widok składa się z trzech sekcji: “Slack”, “Discord” oraz “Email”. W każdej z nich znajduje się lista dodanych kanałów. Poszczególne pola w wierszach zawierają:

- nazwę kanału, która umożliwia łatwe rozróżnienie kanałów,
- webhook powiązany z kanałem (lub adres e-mail w przypadku wiadomości e-mail),
- datę dodania kanału,
- datę ostatniej modyfikacji.

Dodatkowo, każdy wiersz umożliwia użytkownikowi edycję, przetestowanie działania kanału lub jego usunięcie, zgodnie z historyjkami użytkownika. Podobnie jak w innych widokach, ikona „plus” umożliwia dodanie nowego kanału powiadomień.

### 7.6 Diagramy procesów {#diagramy-procesów}

#### 7.6.1 Diagram procesu generowania raportu

<figure>
    <img src="/sequence-diagrams/reports-generation-sequence.svg">
    <figcaption> Widok główny na urządzeniach mobilnych [źródło opracowanie własne]</figcaption>
</figure>

Generowanie raportu jest procesem, który odbywa się na żądanie użytkownika i zakłada przejście przez Aplikację webową (Web Client), Management Service oraz Reports Service.

W ramach tego procesu użytkownik zgłasza żądanie, wykorzystując aplikację webową jako klienta, definiując w niej czy raport powinien być cykliczny (wykonywać się automatycznie co jakiś czas) czy wykonany jednorazowo. Jeżeli raport jest przygotowany na działanie cykliczne, zostanie on wygenerowany co skonfigurowany przez użytkownika okres.. W przeciwnym przypadku, żądanie jest wysyłane natychmiast.

Żądania raportów są przekazywane do przeznaczonych dla tych komunikatów brokerów. Reports Service pobiera żądanie z brokera i rozpoczyna tworzenie raportu z wykorzystaniem zewnętrznego modelu językowego (OpenAI model). Po otrzymaniu rezultatów publikuje raport na odpowiednim brokerze, lub publikuje błąd w przypadku błędu podczas generacji.

Management service pobiera taką odpowiedź z brokera i aktualizuje rekord w bazie odpowiadający raportowi oczekującemu na generowanie.

#### 7.6.2 Diagram procesu zbierania oraz emitowania zmian w metadanych

Metadane są zbierane z klastra Kubernetes przez Agenta, który następnie wysyła je do Metadata Service. Metadata Service zapisuje otrzymane dane w bazie, po czym cyklicznie sprawdza, czy najnowszy stan metadanych aplikacji, hostów lub klastrów uległ zmianie. Jeśli zmiana zostanie wykryta, generowany jest nowy stan, który przesyłany jest bezpośrednio do Management Service.

<figure>
    <img src="/sequence-diagrams/metadata-collection-sequence-diagram.svg">
    <figcaption> Diagram sekwencji zbierania oraz emitowania zmian w metadanych [źródło opracowanie własne]</figcaption>
</figure>

# 8. Implementacja {#implementacja}

## 8.1 Środowisko pracy {#środowisko-pracy}

### 8.1.1 Uruchamiania projektu lokalnie

Ze względu na nature systemu Magpie Monitor i ilość użytych języków i technologii (**18** kontenerów z serwisami i narzędziami), kluczowe było aby uprościć i ustandaryzować proces uruchamiania wszystkich narzędzi oraz mikroserwisów na systemach deweloperskich oraz instancji produkcyjnej.

Wszystkie narzędzia i podsystemy zostały przygotowane w formie kontenerów, które można uruchomić jednym poleceniem `make`.

<figure>
    <img src="/srodowisko-pracy/srodowisko-makefile.png">
    <figcaption> Fragment pliku Makefile [źródło opracowanie własne]</figcaption>
</figure>

### 8.1.2 Automatyczne wykrywanie zmian z docker watch

Dodatkowo konfiguracja została przygotowana aby działać z kontenerami działającymi z takimi samymi obrazami jak te używane w ostatecznym wdrożeniu dzięki funkcjonalności `watch` oferowaną przez interfejs docker-compose, które automatycznie buduje nowy obraz, jeżeli wykryje zmiany w odpowiadającej części kodu. Taka konfiguracja pozwoliła jednorazowo uruchochomić polecenie, które następnie automatycznie wykrywa zmiany z jakiegolwiek z serwisów i odświeża kontener z nowym obrazem. Taka konfiguracja znacząco ułatwiła rozwijanie systemu składającego się z wielu mikroserwisów.

<figure>
    <img src="/srodowisko-pracy/srodowisko-pracy-reports-service.png">
    <figcaption> Przykład konfiguracji Reports Service w docker-compose [źródło opracowanie własne]</figcaption>
</figure>

### 8.1.3 Bogaty zestaw narzędzi do administracji i rozwiązywania problemów

Aby ułatwić pracę z wieloma bazami danych oraz brokerów wiadomości, dodano do konfiguracji docker-compose narzędzia pozwalające administracje tymi systemami z panelu przeglądarki. Wszystkie poniższe narzędzia są częścią polecenia `make watch` dzięki któremu mogą być one uruchomione na systemie dewelopera beż żadnych dodatkowych kroków lub zależności.

Niski koszt tworzenia środowisk z pełnym zestawem narzędzi sprawił, że stały się one bardzo przystępne dla każdego członka zespołu.

Dodatkowo, narzędzia te zostały również udostępnione na publicznej instancji na której uruchomiony był również system Magpie Monitor.

To pozwoliło na administracje i weryfikacje działania wszystkich systemów bez dostępu do samego serwera.

#### Mongo Express

Narzędzie do administracji baz MongoDB. W ramach realizacji Magpie Monitor, narzędzie to było wykorzystywane do administracji baz z raportami z Reports Service oraz Management Service, oraz administracji baz z metadanymi z Cluster Metadata Service oraz Management Service.

<figure>
    <img src="/srodowisko-pracy/srodowisko-pracy-mongo-express.png">
    <figcaption> Interfejs Mongo Express [źródło opracowanie własne]</figcaption>
</figure>

#### Kafka UI

Narzędzie do administracji brokera Apache Kafka. Pozwalało na rozwiązywanie problemów i testowanie funkcjonalności związanych z komunikacją z kolejką. System Magpie Monitor został oparty na komunikacji mikroserwisów za pośrednictwem Apache Kafka, z związku z tym, narzędzie to było kluczowe podczas rozwoju projektu.

<figure>
    <img src="/srodowisko-pracy/srodowisko-pracy-kafka-ui.png">
    <figcaption> Interfejs Kafka UI [źródło opracowanie własne]</figcaption>
</figure>

#### Kibana

Kibana została wykorzystana jako panel administracyjny do bazy ElasticSearch. ElasticSearch był wykorzystany jako narzędzie do przechowywania logów z aplikacji i hostów, w związku z tym, było kluczowe aby móc zwizualizować ilość logów, ich treść oraz indeksy, w ramach których były zapisane, w wygodny i przejrzysty sposób.

<figure>
    <img src="/srodowisko-pracy/srodowisko-pracy-kibana.png">
    <figcaption> Interfejs Kibany [źródło opracowanie własne]</figcaption>
</figure>

#### PG Admin

Panel administracyjny do baz relacyjnych PostgreSQL. Wykorzystywany głównie do weryfikacji funkcjonalności związanych z zarządzaniem użytkownikami oraz konfiguracji raportów w Management Service.

<figure>
    <img src="/srodowisko-pracy/srodowisko-pracy-pgadmin.png">
    <figcaption> Interfejs Pg Admin [źródło opracowanie własne]</figcaption>
</figure>

### 8.1.4 Nacisk na konfigurowalność środowiska uruchomieniowego

Podczas tworzenia wszystkich serwisów, nałożono szczególny nacisk na konfigurowalność wszystkich serwisów wykorzystując zmienne środowiskowe.
Było to istotne aby system był elastyczny na modyfikacje zewnętrznych zależności lub serwisów. Magpie Monitor posiada ponad **120** zmiennych środowiskowych odpowiadających za konfiguracje systemu. Odpowiadają one przykładowo za aspekty takie jak adresy baz danych czy brokerów wiadomości lub konfiguracje związane z przetwarzaniem i paczkowaniem logów.

Pozwoliło to na dostrajanie systemu bez potrzeby modifikacji kodu.

<figure>
    <img src="/srodowisko-pracy/srodowisko-pracy-envs.png">
    <figcaption> Fragment pliku z przykładowymi zmiennymi środowiskowymi [źródło opracowanie własne]</figcaption>
</figure>

### 8.1.4 Automatyczna integracja i wdrażanie od początku rozwoju projektu

Wiele aspektów Magpie Monitor mogło być przetestowanych tylko w systemie imitującym środowisko produkcyjne. W związku z tym kluczowe było aby jak najwcześniej uruchomić system na instancji wystawionej do internetu. Dzięki temu możliwe było zbieranie logów przez całą dobę oraz generowanie raportów w realistycznych warunkach.

Aby usprawnić ten proces, zaimplementowane zostało automatyczne wdrażanie i integracja.

Podczas tworzenia zmian do głównej gałęzi system wersjonowania uruchamiany był proces wdrożeniowy za pośrednictwem **Github Actions**.

W ramach tego procesu budowane były obrazy wszystkich serwisów na podstawie kodu dostępnego w systemie wersjonowania a następnie były one przesyłane do rejestru obrazów.

Po przesłaniu obrazów do rejestru, zmiany były wykrywane przez narzędzie **Watchtower**, które pobierało nowy obraz i restartowało kontenery, które go wykorzystywały.

<figure>
    <img src="/srodowisko-pracy/srodowisko-pracy-watchtower.png">
    <figcaption> Konfiguracja Watchtower [źródło opracowanie własne]</figcaption>
</figure>

### 8.1.5 Wykorzystanie linterów oraz narzędzi do formatowania

W celu utrzymania wspólnego stylu oraz unikania powszechnych błędów związanych ze stylem lub formatowaniem, wymagane było aby każdy członek zespołu korzystał z narzędzi takich jak **ESLint**, **Prettier**, **Stylelint** oraz **go-fmt**. Wykorzystanie zewnętrznych narzędzi pozwoliło uniezależnić się od zintegrowanych środowisk deweloperskich, i zostawić ich wybór dla każdego dewelopera indywidualnie.

### 8.1.6 Inne użyte narzędzia deweloperskie

Podczas rozwoju oprogramowania użyto wielu środowisk deweloperskich w zależności od dewelopra. Od **NeoVim'a**, przez **VSCode**, po **Jetbrains IDEA**. Dzięki wykorzystanie zewnętrznych narzędzi deweloperzy nie byli związani z konkretnym środowiskiem i mogli rozwijać oprogramowanie w takim, który im osobiście najbardziej odpowiadał.

Do wygenerowania interfejsów między mikroserwisami został wykorzystany schemat **AsyncAPI**.

## 8.2 Struktura plików projektu {#struktura-plikow-projektu}

### 8.2.1 Motywacja za wyborem monolitycznego repozytorium

Projekt został stworzony w ramach **monolitycznego repozytorium**.
Podejście to zostało wybrane ze względu na przystępność integracji kolejnych narzędzi i serwisów.
Sprawiło to również, że rozwijanie jakiejkolwiej części projektu stało się przystępniejsze dla każdego dewelopera, ponieważ każdy miał już cały kod i wszystkie narzędzia dostępne w swoim systemie.

<figure>
    <img src="/struktura-plikow-projektu.png">
    <figcaption> Struktura plików projektu [źródło opracowanie własne]</figcaption>
</figure>

### 8.2.2 Struktura repozytorium

- **/.github**: Pliki konfiguracyjne dla Github Actions
- **/agent**: Kod agenta zbierającego logi i metadane z aplikacji i hostów
- **/client**: Kod odpowiedzialny za klienta webowego
- **/docs/diagrams**: Dokumentacja i diagramy projektu
- **/go**: Kod mikroserwisów napisanych w Go: **Metadata Service**, **Reports Service** oraz **Logs Ingestion Service**
- **/management-service**: Kod odpowiedzialny za **Management Service**
- **/nginx-proxy**: Konfiguracja narzędzia **nginx-proxy** do automatycznego proxowania ruchu do serwisów i narzędzi na podstawie domeny oraz terminacji SSL
- **/watchtower**: Konfiguracja narzędzia **Watchtower** od automatycznego wdrażania serwisów
- **/.env.sample**: Plik z przykładowymi wartościami wszystkich zmiennych środowiskowych użytych w projekcie
- **/Makefile**: Plik z poleceniami Make dla uproszczenia wykonywania najpowszechniejszych poleceń podczas rozwoju projektu
- **/docker-compose.yml**: Plik docker-compose zawierający serwisy i narzędzia konieczne do działania Magpie Monitor
- **/docker-compose.es.yml**: Plik docker-compose zawierający konfiguracje ElasticSearch i automatyczne generowanie dla niego certyfikatów kryptograficznych
- **/docker-compose.dev.yml** Plik docker-compose nadpisujący wybrane ustawienia z /docker-compose.yml. Pozwala na efektywniejszy rozwój projektu lokalnie.

### 8.2.3 Struktura plików w katalogu /go

W ramach tego katalogu umieszczono kod wszystkich mikroserwisów napisanych w Go.

Umieszczenie plików z zależnościami (go.mod oraz go.sum) w korzeniu folderu /go pozwoliło na reużywanie zależności oraz wspólnych pakietów z /pkg przez każdy z mikroserwisów. Dzięki temu integracja nowych serwisów w Go była drastycznie przyspieszona.

<figure>
    <img src="/struktura-projektu-go.png">
    <figcaption> Struktura plików w /go [źródło opracowanie własne]</figcaption>
</figure>

- **/go/docker**: Zawiera pliki Dockerfile dla każdego z serwisów
- **/go/pkg**: Zawiera pakiety wspólne dla wszystkich mikroserwisów
- **/go/services**: Każdy z podfolderów ma nazwę odpowiadającą nazwie mikroserwisu, którego kod on posiada

## 8.3 Struktura plików w aplikacji “Agenta” {#struktura-plików-w-aplikacji-“agenta”}

Agent składa się z kilku mniejszych pakietów, z których większość jest współdzielona zarówno przez Pod Agenta jak i Node Agenta.

<figure>
    <img src="/agent/agent-struktura-plikow.png">
    <figcaption> Struktura plików Agenta [źródło opracowanie własne]</figcaption>
</figure>

Wybrane pakiety:

- **agent/app/internal/agent/node** - przechowuje kod odpowiadający za logikę zbierania oraz przesyłania logów z hostów
- **agent/app/internal/agent/pod** - przechowuje kod odpowiadający za logikę zbierania oraz przesyłania logów z aplikacji
- **agent/app/internal/broker** - odpowiada za połączenie z brokerem Kafki
- **agent/app/internal/config** - przechowuje opcje konfiguracyjne Agenta
- **agent/app/internal/database** - odpowiada za połączenie z bazą danych Redis
- **agent/pkg/envs** - udostępnia podstawowe operacje na zmiennych środowiskowych
- **agent/pkg/kubernetes** - udostępnia klienta API Kubernetes
- **agent/pkg/tests** - zawiera implementacje interfejsów używane w testach

W części kodu agenta znajduje się również plik budujący obraz kontenera **Dockerfile** oraz **go.mod** i **go.sum**, które odpowiadają za zarządzanie zależnościami w projekcie języka Golang.

<figure>
    <img src="/agent/agent-struktura-plikow-helm.png">
    <figcaption> Struktura plików paczki wdrożeniowej Helm Agenta [źródło opracowanie własne]</figcaption>
</figure>

Dodatkowo, agent posiada folder chart, który przechowuje Helm Chart, czyli paczkę pozwalającą na łatwe wdrożenie agenta na klastrze Kubernetes. W folderze chart znajduje się również folder scripts, który udostępnia zbiór skryptów przydatnych do testowania agenta.

## 8.4 Struktura plików w serwisie “Logs Ingestion Service” {#struktura-plików-w-serwisie-“ingestion-service”}

Logs Ingestion Service został stworzony zgodnie ze standardami języka go związanych z strukturą projektów. Jednocześnie serwis jest na tyle mały, że nie było potrzeby na nadmierne pakietowanie.

<figure>
    <img src="/logs-ingestion/logs-ingestion-struktura-plikow.png">
    <figcaption> Struktura plików Logs Ingestion Service [źródło opracowanie własne]</figcaption>
</figure>

- **/go/services/logs_ingestion/cmd**: Zawiera pakiet main i tym samym punkt startowy uruchomający serwis
- **/go/services/logs_ingestion/pkg/config**: Zawiera konfiguracje wstrzykiwanych zależności do aplikacji i testów z wykorzystaniem Go fx.
- **/go/services/logs_ingestion/pkg/logs_stream**: Zawiera pakiet logsstream, który jest odpowiedzialny za nasłuchiwanie za logami aplikacji i hostów z brokera wiadomości

## 8.5 Struktura plików w serwisie “Report Service” {#struktura-plików-w-serwisie-“report-service”}

Raports Service został stworzony zgodnie ze standardami języka Go, stąd obecność folderów takich jak /cmd, /internal czy /pkg.

<figure>
    <img src="/reports/raports-struktura-plikow.png">
    <figcaption> Struktura plików Reports Service [źródło opracowanie własne]</figcaption>
</figure>

- **/go/services/reports/cmd**: Zawiera pakiet main i tym samym punkt startowy uruchomający serwis
- **/go/services/reports/api**: Zawiera schematy AsyncAPI do wizualizacji interfejsów serwisu
- **/go/services/reports/internal**: Zawiera kod, który nie może być wykorzystywany jako zależności poza tym katalogiem
  - **/go/services/reports/internal/brokers**: Zawiera konfiguracje wykorzystywanych w serwisie klientów brokerów
  - **/go/services/reports/internal/database**: Zawiera konfiguracje klientów bazy danych wykorzystywanych w serwisie
  - **/go/services/reports/internal/handlers**: Zawiera kod odpowiedzialny za obsługe żądań z zewnętrz
  - **/go/services/reports/internal/services**: Zawiera kod odpowiedzialny za wysokopoziomową logikę biznesową serwisu
- **/go/services/logs_ingestion/pkg**: Zawiera kod, który może być wykorzystywany jako zależności dla innych pakietów
  - **/go/services/logs_ingestion/pkg/config**: Zawiera konfiguracje wstrzykiwanych zależności do aplikacji i testów z wykorzystaniem Go fx.
  - **/go/services/logs_ingestion/pkg/filter**: Zawiera kod odpowiedzialny za filtrowanie logów podczas generowania raportu
  - **/go/services/logs_ingestion/pkg/incident_correlation**: Zawiera kod odpowiedzialny za znajdowanie podobieńst między incydentami i ich scalanie
  - **/go/services/logs_ingestion/pkg/insights**: Zawiera kod odpowiedzialny za generowanie **obserwacji** z logów
  - **/go/services/logs_ingestion/pkg/openai**: Zawiera kod odpowiedzialny za obsługę interfejsów OpenAI
  - **/go/services/logs_ingestion/pkg/repositories**: Zawiera kod odpowiedzialny za przechowywanie trwałych danych w serwisie
  - **/go/services/logs_ingestion/pkg/scheduled_jobs**: Zawiera kod odpowiedzialny za kolejnkowania zadań do modelu językowego

## 8.6 Struktura plików w serwisie “Metadata Service” {#struktura-plików-w-serwisie-“metadata-service”}

<figure>
    <img src="/metadata-service/metadata-service-struktura-plikow.png">
    <figcaption> Struktura plików w Metadata Service [źródło opracowanie własne]</figcaption>
</figure>

Wybrane pakiety Metadata Service odpowiadają za:

- **cluster_metadata/internal/database** - odpowiada za połączenie z bazą danych MongoDB
- **cluster_metadata/internal/handlers** \- wystawia endpoint HTTP odpowiedzialny za kontrolę stanu zdrowia aplikacji
- **cluster_metadata/pkg/repositories** - przechowuje repozytorium danych MongoDB
- **cluster_metadata/pkg/services** - zawiera serwisy odpowiadające za część biznesową aplikacji tj. zbieranie oraz generowanie zagregowanych metadanych, a także emitowanie wydarzeń sygnalizujących zmianę metadanych

## 8.7 Struktura plików w aplikacji “Management Service” {#struktura-plików-w-aplikacji-“management-service”}

<figure>
    <img src="/management-service/management-service-struktura-plikow.png">
    <figcaption> Struktura plików w Management Service [źródło opracowanie własne]</figcaption>
</figure>

Management Service został podzielony na moduły wedle realizowanej przez nie domeny biznesowej, wybrane moduły:

- **auth** - odpowiada za uwierzytelnianie i autoryzację użytkownika, korzystając przy tym z zewnętrznych dostawców uwierzytelniania
- **cluster** - zawiera logikę pozwalającą na odczytywanie klastrów należących do użytkownika oraz przechowywanie i edycję ich konfiguracji
- **metadata** - odbiera, przechowuje oraz udostępnia metadane o klastrach, aplikacjach oraz hostach
- **notifications** - odpowiada za wysyłanie powiadomień oraz konfigurację ich odbiorników
- **reports** - udostępnia możliwość generowania Raportu, przechowuje oraz udostępnia wygenerowane raporty
- **security** - zawiera konfigurację ustawień bezpieczeństwa aplikacji
- **user** - odpowiada za zapisywanie, modyfikowanie oraz tworzenie użytkowników aplikacji
- **utils** - zawiera klasy pomocnicze, agregujące często spotykane w aplikacji funkcje
- **resources** - rzechowuje konfigurację uruchomieniową aplikacji oraz szablony stylizacji wysyłanych powiadomień
- **test** - zawiera testy jednostkowe oraz integracyjne

Plik **Dockerfile** zawiera konfigurację budowania obrazu, a pliki **mvnw** i **pom.xml** należą do narzędzia Maven, które zarządza zależnościami w projekcie i buduje plik JAR aplikacji.

## 8.8 Struktura plików w aplikacji klienckiej {#struktura-plików-w-aplikacji-klienckiej}

Struktura aplikacji klienckiej została zaprojektowana zgodnie z najlepszymi praktykami React oraz Typescript.

<figure>
    <img src="/client/client-struktura-plikow.png">
    <figcaption> Struktura plików w aplikacji klienckiej [źródło opracowanie własne]</figcaption>
</figure>

- **/client/index.html**: Główny plik index.html, który jest zwracany przez serwer webowy
- **/client/.nginx**: Katalog zawierający konfiguracje Nginx, który służy jako serwer webowy udostępniający pliki statyczne klienta
- **/client/.eslintrc**: Plik konfiguracyjny narzędzia ESLint
- **/client/.stylelintignore**: Plik konfiguracyjny narzędzia Stylelint
- **/client/.stylelintrc**: Plik konfiguracyjny narzędzia Stylelint
- **/client/Dockerfile**: Plik Dockerfile definiujący obraz kontenera z serwerem klienta webowego
- **/client/package-lock.json**: Plik definiujący "zamrożone" zależności projektu.
- **/client/package.json**: Plik definiujący zależności projektu
- **/client/tsconfig.json**: Plik konfiguracyjny kompilatora Typescript
- **/client/tsconfig.node.json**: Plik konfiguracyjny kompilatora Typescript
- **/client/vite.config.json**: Plik konfiguracyjny bundlera Vite
- **/client/public**: Katalog z plikami dostępnymi do pobrania przez skrypt z przeglądarki użytkownika
- **/client/src**: Katalog z plikami skryptów
  - **/client/src/api**: Katalog z logiką odpowiedzialną za komunikacje z zewnętrznymi interfejsami
  - **/client/src/assets**: Katalog z grafikami używanymi z aplikacji
  - **/client/src/components**: Katalog z reużywalnymi komponentami
  - **/client/src/global**: Katalog z globalnymi zmiennymi
  - **/client/src/hooks**: Katalog z customowymi hookami
  - **/client/src/lib**: Katalog z reużywalną logiką
  - **/client/src/links**: Katalog z linkami wyświetlanymi użytkownikowi w aplikacji
  - **/client/src/messages**: Katalog z komunikatami prezentowanymi użytkownikowi
  - **/client/src/pages**: Katalog z komponentami stron
  - **/client/src/providers**: Katalog z definicją kontekstów i ich dostawców
  - **/client/src/types**: Katalog z globalnymi typami używanymi w aplikacji
  - **/client/src/index.scss**: Katalog z globalnymi stylami
  - **/client/src/main.tsx**: Plik z punktem wejściowym skrytpu
  - **/client/src/variables.scss**: Katalog z globalnymi zmiennymi SCSS.

## 8.9 Uwierzytelnienie użytkownika {#uwierzytelnienie-użytkownika}

Do uwierzytelniania użytkowników w aplikacji wykorzystywany jest dostawca OAuth2 od Google. Proces weryfikacji dostępu przebiega w kilku etapach:

#### Proces dla użytkownika logującego się do aplikacji:

1. Użytkownik otwiera stronę główną aplikacji.
2. Wybierany jest sposób logowania, w tym przypadku logowanie przy użyciu Google.
3. Następuje przekierowanie na ekran logowania dostarczony przez Google.
4. Użytkownik wprowadza dane logowania do swojego konta Google.
5. Po zalogowaniu użytkownik wyraża zgodę na udostępnienie aplikacji danych osobowych, takich jak imię, nazwisko i adres e-mail.
6. Po pomyślnym uwierzytelnieniu użytkownik zostaje przekierowany z powrotem do aplikacji.
7. Mikroserwis `management-service` otrzymuje trzy tokeny: token autoryzacyjny, **ID token**, oraz **refresh token**. W dalszym procesie wykorzystywane są dwa z nich:
   - **ID token** – ważny przez jedną godzinę i dołączany do każdego zapytania HTTP do `management-service`.
   - **Refresh token** – ważny bezterminowo i używany wyłącznie do odświeżania ID tokenu. Tokeny są przechowywane w przeglądarce użytkownika jako ciasteczka: `authToken` (ID token) oraz `refreshToken` (refresh token). Taki podział zapewnia ograniczenie ryzyka wycieku refresh tokena, który umożliwia generowanie nowych tokenów uwierzytelniających.

#### Obsługa zapytań HTTP:

8. Każde zapytanie przychodzące z frontendu do `management-service` jest filtrowane. Wyjątek stanowią następujące zasoby:
   - `/public/api/*`
   - `/v3/api-docs/*`
   - `/swagger-ui/*`
   - `/api/v1/auth/*`
   - `/login/oauth2/code/*`
   - `/error`
     Dostęp do tych zasobów nie wymaga uwierzytelnienia. Wszystkie inne endpointy wymagają dołączenia ważnego ID tokenu w nagłówku zapytania.

#### Odświeżanie tokenu uwierzytelnienia:

9. W przypadku wygaśnięcia ID tokenu możliwe jest jego odświeżenie przez wywołanie endpointu:
   ```
   https://management-service.rolo-labs.xyz/api/v1/auth/refresh-token
   ```
   W zapytaniu należy przekazać refresh token. W odpowiedzi użytkownik otrzymuje nowy ID token. Proces ten jest zautomatyzowany – frontend regularnie sprawdza długość ważności bieżącego ID tokenu, korzystając z endpointu:
   ```
   https://management-service.rolo-labs.xyz/api/v1/auth/auth-token/validation-time
   ```
   Informacja o czasie ważności tokenu jest odczytywana z ładunku ID tokenu, który ma format JWT (JSON Web Token). Jeśli token wygaśnie w ciągu najbliższych 30 minut, aplikacja automatycznie wysyła żądanie o nowy ID token.

#### Wylogowanie użytkownika:

Użytkownik może wylogować się z aplikacji w dowolnym momencie. Po kliknięciu przycisku wylogowania oba ciasteczka (`authToken` oraz `refresh`) są trwale usuwane z przeglądarki.

Taki model uwierzytelnienia zapewnia łatwość obsługi dla użytkownika oraz zgodność z najlepszymi praktykami bezpieczeństwa, minimalizując ryzyko nieautoryzowanego dostępu do aplikacji.

## 8.10. Planowanie raportów (scheduling raportów, management service) {#planowanie-raportów-(scheduling-raportów,-management-service)}

Management Service udostępnia użytkownikowi funkcję konfiguracji raportów cyklicznych, które będą generowane co zdefiniowany przez użytkownika okres.

<figure>
    <img src="/management-service/management-service-generate-reports.png">
    <figcaption>Cykliczne generowanie zaplanowanych raportów</figcaption>
</figure>

Cyklicznie uruchamiany komponent sprawdza, czy

<figure>
    <img src="/management-service/management-service-process-schedule.png">
    <figcaption>Przetwarzanie zaplanowanego raportu</figcaption>
</figure>

<figure>
    <img src="/management-service/management-service-create-report.png">
    <figcaption>Generowanie raportu</figcaption>
</figure>


<figure>
    <img src="/management-service/management-service-publish-report-requested.png">
    <figcaption>Publikowanie wydarzenia tworzącego raport</figcaption>
</figure>

## 8.11 Zbieranie logów (agent) {#zbieranie-logów-(agent)}

## 8.11.1 Zbieranie logów z hostów

Agent zbiera logi z hostów, obserwując pliki które użytkownik skonfiguruje podczas wdrożenia.

<figure>
    <img src="/agent/agent-configuration.png">
    <figcaption> Konfiguracja Agenta [źródło opracowanie własne]</figcaption>
</figure>

Konfiguracja plików logów odbywa się w pliku values.yaml paczki wdrożeniowej Helm Chart.

<figure>
    <img src="/agent/agent-watch-files.png">
    <figcaption> Obserwowanie plików przez Agenta [źródło opracowanie własne]</figcaption>
</figure>

Podczas uruchomienia, node Agent obserwuje wszystkie skonfigurowane pliki.

<figure>
    <img src="/agent/agent-watch-file.png">
    <figcaption> Obserwowanie pliku przez Agenta [źródło opracowanie własne]</figcaption>
</figure>

Obserwowanie pliku polega na cyklicznym sprawdzeniu jego rozmiaru i porównaniu z rozmiarem poprzedniego odczytu, który zapisany jest w bazie danych Redis. Jeśli rozmiar zmienił się względem poprzedniego odczytu, agent odczytuje nieodczytane wcześniej dane.

<figure>
    <img src="/agent/agent-split-into-packets.png">
    <figcaption> Podział logów na pakiety [źródło opracowanie własne]</figcaption>
</figure>

<figure>
    <img src="/agent/agent-node-chunk.png">
    <figcaption> Pakiet danych [źródło opracowanie własne]</figcaption>
</figure>

Dane dzielone są na mniejsze pakiety _Chunk_, które przesyłane są do brokera Kafki.

## 8.11.2 Zbieranie metadanych z hostów

<figure>
    <img src="/agent/agent-gather-node-metadata.png">
    <figcaption> Zbieranie metadanych o hostach [źródło opracowanie własne]</figcaption>
</figure>

<figure>
    <img src="/agent/agent-node-metadata.png">
    <figcaption> Pakiet metadanych [źródło opracowanie własne]</figcaption>
</figure>

Agent cyklicznie odczytuje oraz przesyła metadane o swoim działaniu, tj. id klastra na którym działa, nazwę hosta w klastrze na którym działa dana replika agenta oraz pliki obserwowane przez agenta. Odczytane metadane są następnie przesyłane do brokera Kafki.

## 8.11.3 Zbieranie logów z aplikacji

Logi aplikacji są zbierane ze wszystkich przestrzeni nazw obecnych w klastrze Kubernetes, poza przestrzeniami wykluczonymi w konfiguracji paczki wdrożeniowej Helm Chart.

<figure>
    <img src="/agent/agent-pod-configuration.png">
    <figcaption> Konfiguracja Pod Agenta [źródło opracowanie własne]</figcaption>
</figure>

Logi zbierane są cyklicznie, co zadany, konfigurowalny okres dla każdej z włączonych przestrzeni nazw.

<figure>
    <img src="/agent/agent-gather-logs.png">
    <figcaption> Zbieranie logów aplikacji [źródło opracowanie własne]</figcaption>
</figure>

W każdej z przestrzeni, agent pobiera logi z obiektów Kubernetes typu Deployment, StatefulSet oraz DaemonSet.

<figure>
    <img src="/agent/agent-fetch-logs-for-namespace.png">
    <figcaption> Zbieranie logów dla przestrzeni nazw [źródło opracowanie własne]</figcaption>
</figure>

<figure>
    <img src="/agent/agent-fetch-deployment-logs-since-time.png">
    <figcaption> Zbieranie logów zasobu Deployment [źródło opracowanie własne]</figcaption>
</figure>

<figure>
    <img src="/agent/agent-fetch-pod-logs-since-time.png">
    <figcaption> Zbieranie logów zasobu Deployment [źródło opracowanie własne]</figcaption>
</figure>

<figure>
    <img src="/agent/agent-get-pod-log-packets.png">
    <figcaption> Podział logów na pakiety [źródło opracowanie własne]</figcaption>
</figure>

<figure>
    <img src="/agent/agent-application-chunk.png">
    <figcaption> Pakiet danych [źródło opracowanie własne]</figcaption>
</figure>

Powyższe obrazki przedstawiają przykład zbierania logów z obiektu typu Deployment. Dla danego Deploymentu pobierane są wszystkie należące do niego pody, z których logi są pobierane oraz dzielone na pakiety Chunk. Pakiety są następnie przesyłane do brokera Kafki. Analogiczny proces zachodzi również dla obiektów typu StatefulSet oraz DaemonSet.

<figure>
    <img src="/agent/agent-kubernetes-client.png">
    <figcaption> Klient API Kubernetes [źródło opracowanie własne]</figcaption>
</figure>

## 8.11.4 Zbieranie metadanych z aplikacji

Agent cyklicznie odczytuje oraz przesyła metadane o aplikacjach aktualnie działających na klastrze. Odczytane metadane są następnie przesyłane do brokera Kafki.

<figure>
    <img src="/agent/agent-gather-cluster-metadata.png">
    <figcaption> Zbieranie metadanych o aplikacjach [źródło opracowanie własne]</figcaption>
</figure>

<figure>
    <img src="/agent/agent-application-metadata.png">
    <figcaption> Pakiet metadanych [źródło opracowanie własne]</figcaption>
</figure>

## 8.12 Zapisywanie logów (ingestion service) {#zapisywanie-logów-(ingestion-service)}

### 8.12.1 Dynamiczne tworzenie indeksów

Wiadomość z logami może posiadać parę poziomów zagnieżdżenie, grupujących logi na podstawie określonych metadanych. W związku z tym kluczowe jest przetworzenie tej wiadomości w taki sposób, aby dokument nie miał zagnieżdżonych pół, ponieważ utrudnia to przeszukiwania i indeksowanie. W ramach takiego procesu, z jednej wiadomości może powstać wiele dokumentów.

<figure>
    <img src="/logs-ingestion/implementation/logs-ingestion-implementation-index-name.png">
    <figcaption> Implementacja Logs Ingestion: Nazwa indeksu  [źródło opracowanie własne]</figcaption>
</figure>

### 8.12.2 Spłaszczanie dokumentów

W ramach jednej wiadomości zawierającej logi sama zawartość może być zagnieżdżona, w związku z tym kluczowe jest przetworzenie tej wiadomości w taki sposób aby dokument nie miał zagnieżdżonych pół, ponieważ utrudnia to przeszukiwania i indeksowanie. W ramach takiego procesu z jednej wiadomości może powstać wiele dokumentów.

<figure>
    <img src="/logs-ingestion/implementation/logs-ingestion-implementation-flatten.png">
    <figcaption> Implementacja Logs Ingestion: Spłaszczanie logów  [źródło opracowanie własne]</figcaption>
</figure>

### 8.12.3 Grupowanie operacji dodawania rekordów

W związku z tym, że w jednej wiadomości może być wiele dokumentów, to aby uniknąć wykonywania wielu zapytań do bazy, wykorzystany został mechanizm Bulk Query, który pozwala wysyłać listę zadań ElasticSearch do wykonania.

<figure>
    <img src="/logs-ingestion/implementation/logs-ingestion-implementation-inserting-logs.png">
    <figcaption> Implementacja Logs Ingestion: Dodawanie logów [źródło opracowanie własne]</figcaption>
</figure>

## 8.13. Generowanie raportów (reports service) {#generowanie-raportów-(reports-service)}

### 8.13.1 **Wczytywanie logów do pamięci**

Istotnym problemem, występującym podczas pracy z dużej ilością danych jest unikanie sytuacji, w których musimy mieć wszystkie dane w pamięci jednocześnie. Ze względu na wysoki rozmiar logów, z których użytkownik może sobie zażyczyć raportu, kluczowe było aby przetwarzać je w paczkach, później nazywanych batchami.

Proces ten rozpoczyna się przez pobieranie logów z instancji ElasticSearch w batchach wykorzystując dostarczone Scroll API. Interfejs ten pozwala na stworzenie zapytania, dla którego zostanie zwrócony identyfikator scrolla, który może służyć do pobierania kolejnych paczek z zapytania, realizując w ten sposób paginacje zapytań.

Wykorzystując ten mechanizm możliwe było stworzenie interfejsu obejmującego pobieranie logów w paczkach z dowolnej bazy danych.

<figure>
    <img src="/reports/implementation/reports-implementation-batched-document-retriever.png">
    <figcaption> Interfejs do paczkowania dokumentów (logów) [źródło opracowanie własne]</figcaption>
</figure>

W ten sposób możliwe jest operowanie wyłącznie na paczkach danych, i zdejmowanie ich z pamięci wraz z rozpoczęciem działania na następnej paczce.

<figure>
    <img src="/reports/implementation/reports-implementation-next-batch.png">
    <figcaption> Reports Service: Zdobywanie następnej paczki logów [źródło opracowanie własne]</figcaption>
</figure>

### 8.13.2 Przygotowywanie logów do generowania raportu

Logi przetwarzane są po jednej paczce w danym momencie (co sprowadza się w obecnej konfiguracji do 10 000 rekordów) Przed przesłaniem żądania, logi są grupowane oraz filtrowane, tak aby zminimalizować koszta wygenerowania raportu i zmniejszenie zbędnego obciążenia.

<figure>
    <img src="/reports/implementation/reports-implementation-baching-logs.png">
    <figcaption> Reports Service: Przetwarzanie logów paczkami [źródło opracowanie własne]</figcaption>
</figure>

### 8.13.3 Tworzenie zapytań do modelu językowego od OpenAI

Główną funkcjonalnością **Reports service** jest przekształcanie logów w raporty, z wykorzystaniem zewnętrznego modelu językowego.

Stąd, bardzo istotny był odpowiedni sposób na komunikację z modelem pozwalający na utrzymanie niskich kosztów, uwzględniając ilość logów, którą model musiałby przetworzyć.

Ze względu na łatwą integrację i wszechstronność wykorzystany został model gpt-4-mini od OpenAI. Model ten został wybrany ze względu na niską cenę zapytań: $1.25 na milion tokenów.

Dodatkowym czynnikiem było wsparcie dla ustrukturyzowanej odpowiedzi. Jest to funkcjonalność oferowana przez modele od OpenAI (gpt-4-mini to najtańszy model wspierający tą funkcjonalność), gwarantująca, że odpowiedź modelu będzie miała formę konkretnego schematu (np. JSON). Było to kluczowe ze względu na identyfikatory logów na podstawie których wykryty został incident. Wymagane było aby znajdowały się one w odpowiedzi od modelu. Ustrukturyzowane odpowiedzi to jedyny pewny sposób na ich otrzymanie w przewidywalnym formacie.

Ze względu na brak oficjalnego klienta do API OpenAI w Go, napisany został własnościowy klient wspierający serializację strukur z Go do schematów definiujących w jaki sposób model może odpowiedzieć. Aby to osiągnąć wykorzystano mechanizmy refleksji dostępne w języku Go.

<figure>
    <img src="/reports/implementation/reports-implementation-structured-output.png">
    <figcaption> Reports Service: Implementacja transformacji do ustrukturyzowanego wyjścia modelu [źródło opracowanie własne]</figcaption>
</figure>

W ten sposób możliwa była automatyczna serializacja struktur obecnych w projekcie do schematów w których odpowiada model

Tak stworzony schemat przesłany razem z zapytaniem do modelu pozwala na zagwarantowanie odpowiedzi spełniającej ten kontrakt.

Dodatkowym kluczowym mechanizmem oferowanym przez OpenAI jest Batch API. Jest to mechanizm pozwalający na obniżenie kosztów zapytań do **50%**, oraz grupowania zapytań w listy, które będą realizowane wspólnie do 24 godzin od momentu zgłoszenia. Jest to mechanizm stworzony dla systemów przetwarzających dużą ilość danych oraz jednocześnie niekrytycznych pod względem czasu wykonania. W naszym przypadku było to idealne rozwiązanie, ponieważ logi pod względem objętości są wyjątkowo ciężkie. Jednocześnie generowanie raportów to coś co jest wykonywanie cyklicznie (zwykle co parę dni), więc nie występuje presja czasu, które by zdyskwalifikowała to rozwiązanie.

### 8.13.3 Obsługa interfejsu Batch API

Tak jak wcześniej wspomniano, zrealizowanie komunikacji z Batch API od OpenAI wymagało implementacji własnościowego klienta, a w ramach niego metod na serializacje wielu zapytań do jednego żądania.

Do modelu przesyłana była zserializowana w formacie **JSONL** lista struktur **BatchFileCompletionRequestEntry**.

Element taki posiada identyfikator przypisywany przez klienta **CustomId**, który pozwala powiązać żądanie z odpowiedzią będącą w formacie listy struktur **BatchFileCompletionResponseEntry**. Jest to konieczne, ponieważ Batch API nie gwarantuje zachowania takiej samej kolejności odpowiedzi w jakiej zostały zakodowane zapytania.

**JSONL** to format przeznaczony do kodowania listy jsonów, gdzie każdy kolejny JSON jest rozdzielony znakiem nowej linii.

Ten format również wymagał napisania własnościowego enkodera wykorzystującego mechanizmy refleksji obecne w Go.

<figure>
    <img src="/reports/implementation/reports-implementation-batch-api-interface.png">
    <figcaption> Reports Service: Interfejs Batch API [źródło opracowanie własne]</figcaption>
</figure>

### 8.13.4 Wykorzystanie FilesAPI razem z BatchAPI

Przed przesłaniem żądania przetworzenia batcha przez BatchAPI, konieczne jest przesłanie treści żądań (w formacie **JSONL**), jako plik do Files API od OpenAI. Po przesłaniu pliku, możliwe jest rozpoczęcie przetwarzania batcha odwołując się do identyfikatora pliku z zapytaniami.

Dodatkowym parametrem wymagającym sprecyzowania jest okno wykonania, czyli maksymalny czas w jakim ma być wykonane dane żądanie. Niestety w momencie realizacji tego projektu jedyną wspieraną wartością dla okna wykonania były 24 godziny.

<figure>
    <img src="/reports/implementation/reports-implementation-batch-api-process.png">
    <figcaption> Reports Service: Wykorzystanie Batch API [źródło opracowanie własne]</figcaption>
</figure>

### 8.13.5 Dzielenie zapytań do modelu na konteksty i batche

Podstawowym problemem podczas interpretacji logów jest ograniczenie wielkości logów przekazanych w ramach jednego kontekstu. Zmniejszenie ilości logów wewnątrz pojedyńczego kontekstu pozwoliła zaobserwować poprawę jakości wykrywanych incydentów i rekomendacji ich rozwiązania.

Dodatkowym twardym ograniczeniem jest maksymalna wielkość kontekstu, która wynosi około 100 000 tokenów, który przekazując logi można bardzo łatwo przekroczyć.

Aby utrzymać wysoką jakość odpowiedzi, zdecydowano się pogrupować logi na podstawie aplikacji / hostów, przez, które zostały wyprodukowane oraz następnie podzielenić je dodatkowo tak aby jedna paczka logów nie przekraczała z góry określonej wielkości, która jest mniejsza od maksymalnej wielkości kontekstu.

Dzięki temu uzyskujemy konteksty, które dotyczą logów z wyłącznie jednej aplikacji / hosta, a jednocześnie są na tyle małe, że model nie zapomina pojedyńczych elementów loga. To pozwala uniknąć halucynacji i niskiej jakości rekomendacji.

<figure>
    <img src="/reports/implementation/reports-implementation-splitting-logs.png">
    <figcaption> Reports Service: Dzielenie logów do wielu zapytań [źródło opracowanie własne]</figcaption>
</figure>

Dodatkowym ograniczeniem jest maksymalny rozmiar Batcha do OpenAI (Batch API), który wynosi 2MB.

Aby to obejść zapytania są dzielone na paczki zapytań nieprzekraczające tej wielkości.

Tak stworzona paczka stanowi jednostkę pracy i jest zapisywana do bazy, gdzie jest ona następnie wykrywana przez **BatchPoller** (komponent serwisu **Reports Service**) i przekształcana na **obserwacje** przy użyciu modelu jęzkykowego od OpenAI.

Abstrakcje nad jednostką pracą stanowią obiekty ScheduledJob, który posiada zserializowane żądanie do modelu językowego.

 <figure>
    <img src="/reports/implementation/reports-implementation-split-completion-requests.png">
    <figcaption> Reports Service: Dzielenie logów do wielu zapytań [źródło opracowanie własne]</figcaption>
</figure>

### 8.13.6 Kolejkowanie zapytań do modelu

Oferowany przez OpenAI Batch API posiada ograniczenie na liczbę obecnie przetwarzanych tokenów (około 2 miliony tokenów). Każdy batch przekraczający ten limit jest automatycznie odrzucany. W związku z tym wymagane było stworzenie rate-limitera sprawdzającego ile tokenów jest obecnie przetwarzanych przed przesłaniem kolejnego batcha.

 <figure>
    <img src="/reports/implementation/reports-implementation-openaijob.png">
    <figcaption> Reports Service: Schemat jednostki przetwarzania przez model od OpenAI [źródło opracowanie własne]</figcaption>
</figure>

W ten sposób powstała kolejka, która posiada obiekty typu OpenAiJob, o określonym statusie:

- **ENQUEUED**: Żądanie oczekuje na przesłanie do Batch API
- **IN_PROGRESS**: Żądanie jest obecnie przetwarzanie przez OpenAI
- **COMPLETED**: Żądanie zostało zakończone i rezultaty są dostępne za pośrednictwem Batch API
- **FAILED**: Żądanie zostało odrzucone przez OpenAI

Kolejka ta przechowywana jest w bazie. Pozwala to na zachowanie stanu oczekujących zadań, nawet po awarii systemu.

Podczas uruchamiania serwisu, uruchamiany jest również wątek, którego odpowiedzialnością jest sprawdzanie oczekujących zadań z bazy, i wykonywanie następnych zadań jeżeli tokeny obecnie wykonywanych żądań i oczekującego żądania są poniżej limitu.

 <figure>
    <img src="/reports/implementation/reports-implementation-splittiong-batches.png">
    <figcaption> Reports Service: Kolejkowanie zapytań do modelu językowego [źródło opracowanie własne]</figcaption>
</figure>

### 8.13.7 **Interfejs generowania raportu z modelem**

Generowanie raportu z użyciem modelu językowego, polega na przekazaniu modelowi logów wraz z informacjami o ich źródle (aplikacja/host i ich metadane), oraz wskazanie na wejściu do modelu polecenia precyzującego kontekts przesyłanych logów oraz polecenia odnoczące się do sposobu ich analizy oraz formatu odpowiedzi.

Dodatkowym parametrem przekazywanym do wejścia modelu jest konfigurowalny **customPrompt**, który pozwala dostosować interpretacje logów do poszczególnych aplikacji lub hostów.

Wymagany schemat odpowiedzi modelu zawiera listę wykrytych incydent, a w każdym z nich, między innymi: jego tytuł, podsumowanie, rekomendację dotyczącą rozwiązania oraz identyfikatory logów, z których dany incydent został wydedukowany.

<figure>
    <img src="/reports/implementation/reports-implementation-generation-prompt.png">
    <figcaption> Reports Service: Wejście do modelu językowego [źródło opracowanie własne]</figcaption>
</figure>

### 8.13.8 **Scalanie incydentów**

Tak jak wcześniej wspomniano ograniczony rozmiar kontekstu modelu językowego powoduje, że wiele incydentów zostanie powielonych ze względu na to, że logi wyprodukowane na skutek tego samego problemu mogą być analizowane w ramach wielu kontekstów (tak aby nie przekroczyć jego maksymalnego rozmiaru).

Aby uniknąć zduplikowanych incydentów (takich, które odnoszą się faktycznie do tego samego problemu, ale zawierają logi należące do różnych kontekstów), konieczne było ich scalanie.

Incydenty odnoszące się do tego samego problemu posiadają podobne tytuły oraz podsumownia, ale przez to, że były one generowane przez model językowe, który działa niederministycznie, to nie są one identyczne.

W związku z tym konieczne było scalanie incydentów na podstawie rozumienia ich treści. Podejściem pozwalajacym na rozumienie tekstu i łączenie incydentów na podstawie ich treści jest wykorzystanie modelu językowego.

Zadanie scalania incydentów polega na przekazaniu do modelu językowego wyłącznie listy incydentów uproszczonych do identyfikatora, tytułu oraz podsumowania. Na podstawie tej treści model grupuje incydenty ze względu na ich podobieństwa.

Przez mały rozmiar przekazywanych parametrów, możliwe jest aby przekazać wszystkie incydenty z danej aplikacji/hosta w ramach pojedycznego kontekstu modelu. To pozwala na uniknięcie pomijania scalania incydentów należących do różnych kontekstów (analogiczny problem do tego który wymuszał scalanie incydentów).

<figure>
    <img src="/reports/reports-service-merged-incidents.png">
    <figcaption> Reports Service: Scalanie incydentów [źródło opracowanie własne]</figcaption>
</figure>

## 8.14 Ustawianie kanałów komunikacji (management service) {#ustawianie-kanałów-komunikacji-(management-service)}

## 8.15 Odczytywanie stanu klastra (metadata service) {#odczytywanie-stanu-klastra-(metadata-service)}

Głównym zadaniem Metadata Service jest odczytywanie stanu klastra, który definiujemy jako:

- zbiór działających na nim aplikacji
- zbiór hostów należących do klastra oraz obserwowanych przez nie plików z logami

Dodatkowo, Metadata Service zajmuje się agregacją odczytanych metadanych.

## 8.15.1 Odbieranie metadanych o aplikacjach oraz hostach

Metadane o aplikacjach oraz hostach są odbierane z brokera Kafki, a następnie zapisywane do bazy danych MongoDB.

<figure>
    <img src="/metadata-service/metadata-application-state.png">
    <figcaption> Odbierane metadane o aplikacjach [źródło opracowanie własne]</figcaption>
</figure>

Przykładowe metadane zbierane o aplikacjach.

<figure>
    <img src="/metadata-service/metadata-consume-application-metadata.png">
    <figcaption> Odbierane metadane o aplikacjach [źródło opracowanie własne]</figcaption>
</figure>

Powyższy obrazek przedstawia proces zbierania metadanych dla aplikacji, proces ten jest analogiczny dla hostów.

### 8.15.2 Agregacja metadanych dla aplikacji

Zapisane do bazy metadane są następnie odczytywane przez proces agregacji. Agregacja jest procesem cyklicznym, który odbywa się co zadany, konfigurowalny przedział czasowy. Celem procesu agregacji jest przeprowadzenie procesu spłaszczenia otrzymanych metadanych dla każdego zarejestrowanego klastra. Spłaszczone metadane są następnie porównywane z ostatnim odczytem, aby w przypadku zmiany wyemitować zdarzenie zawierające najnowszy stan.

<figure>
    <img src="/metadata-service/metadata-poll-for-application-state-change.png">
    <figcaption> Cykliczne sprawdzanie zmian stanu [źródło opracowanie własne]</figcaption>
</figure>

Powyższy obrazek przedstawia cyklicznie uruchamiany proces agregacji metadanych dla aplikacji.

<figure>
    <img src="/metadata-service/metadata-update-application-metadata-state-for-cluster.png">
    <figcaption> Aktualizacja metadanych aplikacji dla klastra [źródło opracowanie własne]</figcaption>
</figure>

Powyższy obrazek przedstawia przebieg procesu agregacji, w którym odczytywany jest ostatni zagregowany stan, który następnie jest porównywany ze stanem aktualnym. W przypadku zmiany, generowany jest zagregowany stan, który zapisywany jest w bazie danych oraz emitowany w postaci wydarzenia do brokera Kafki.

<figure>
    <img src="/metadata-service/metadata-generate-aggregated-application-state-for-cluster.png">
    <figcaption> Generowanie zagregowanych metadanych aplikacji dla klastra [źródło opracowanie własne]</figcaption>
</figure>

Powyższy obrazek przedstawia proces generowania nowego, zagregowanego stanu w przypadku zmiany względem poprzedniego.

<figure>
    <img src="/metadata-service/metadata-aggregated-application-metadata.png">
    <figcaption> Zagregowane metadane aplikacji [źródło opracowanie własne]</figcaption>
</figure>

<figure>
    <img src="/metadata-service/metadata-publish-application-metadata-updated-event.png">
    <figcaption> Publikowanie wydarzenia ze zmienionym stanem [źródło opracowanie własne]</figcaption>
</figure>

Wynik procesu agregacji metadanych jest emitowany w postaci wydarzenia do brokera Kafki.

Analogiczny proces przeprowadzany jest dla agregacji metadanych hostów.

### 8.15.3 Agregacja metadanych dla klastrów

Dla agregacji metadanych klastrów, proces nieco się różni. Unikalny zbiór klastrów jest odczytywany na podstawie ich identyfikatorów, zawartych w metadanych aplikacji oraz hostów, przesłanych do Metadata Service w zadanym przedziale czasowym.

<figure>
    <img src="/metadata-service/metadata-get-unique-cluster-ids-for-period.png">
    <figcaption> Pobranie identyfikatorów wszystkich klastrów [źródło opracowanie własne]</figcaption>
</figure>

Powyższy obrazek przedstawia pobranie unikalnego zbioru klastrów z danego okresu. Proces bazuje na zebranych metadanych aplikacji oraz hostów, a jego wynik wykorzystywany jest do dalszego porównania, czy zbiór owy uległ zmianie względem ostatniego odczytu.

<figure>
    <img src="/metadata-service/metadata-compare-cluster-states.png">
    <figcaption> Porównanie poprzedniego oraz aktualnego stanu klastra [źródło opracowanie własne]</figcaption>
</figure>

Przykład porównania, na podstawie którego w przypadku zmiany generowany jest nowy zagregowany stan.

<figure>
    <img src="/metadata-service/metadata-aggregated-cluster-metadata.png">
    <figcaption> Zagregowane metadane klastrów [źródło opracowanie własne]</figcaption>
</figure>

<figure>
    <img src="/metadata-service/metadata-publish-cluster-metadata-updated-event.png">
    <figcaption> Publikowanie wydarzenia ze zmienionym stanem [źródło opracowanie własne]</figcaption>
</figure>

Wygenerowany stan jest następnie emitowany w postaci wydarzenia do brokera Kafki.

## 8.16 Zabezpieczenia aplikacji (management service) {#zabezpieczenia-aplikacji-(management-service)}

Aplikacja wykorzystuje protokół OAuth2 do uwierzytelniania użytkowników. Wybrano dostawcę Google, który odpowiada za autoryzację i generowanie tokenów uwierzytelniających. Po zakończonym procesie autoryzacji użytkownik otrzymuje dwa tokeny:

**ID Token (authToken)** - jest to token w formacie JWT, który zawiera podstawowe informacje o użytkowniku oraz potwierdza jego tożsamość. Ten token jest przechowywany w ciasteczku i wysyłany z każdym zapytaniem do serwera aplikacji w celu uwierzytelnienia użytkownika i udostępnienia żądanych zasobów. Czas ważności tego tokena to 1h.

**Refresh Token (refreshToken)** - jest przechowywany w ciasteczku i wykorzystywany wyłącznie w procesie odświeżania sesji. Umożliwia on generowanie nowego ID Tokenu, gdy poprzedni wygasa, bez konieczności ponownego logowania użytkownika. Token ten jest ważny bezterminowo.

Aplikacja nie implementuje podziału użytkowników na role. Wszyscy użytkownicy, którzy przejdą proces uwierzytelnienia, mają jednakowy poziom dostępu do zasobów i funkcjonalności aplikacji.

Proces uwierzytelniania przebiega w następujący sposób:

1. Użytkownik przekierowywany jest na stronę logowania Google, gdzie podaje swoje dane logowania.
2. Po pomyślnym uwierzytelnieniu Google zwraca ID Token i Refresh Token.
3. ID Token zapisywany jest w ciasteczku jako `authToken` i wysyłany w nagłówku każdego zapytania HTTP, aby serwer aplikacji mógł uwierzytelnić użytkownika.
4. Refresh Token przechowywany jest w ciasteczku jako `refreshToken` i wykorzystywany wyłącznie przy odświeżaniu tokenu (np. po jego wygaśnięciu).

Mechanizm ten zapewnia prostotę implementacji uwierzytelniania, przy jednoczesnym wykorzystaniu standardowych rozwiązań OAuth2 i Google Identity Platform.

## 9. Testy produktu programowego/Wyniki i analiza badań {#testy-produktu-programowego/wyniki-i-analiza-badań}

### 9.1 Testy Reports Service {#testy-reports-service}

#### 9.1.1 Testy jednostkowe {#testy-jednostkowe-reports-service}

**Id**: TC1
**Tytuł**: Wydobywanie indexów ElasticSearch na podstawie daty, źródła logów oraz identyfikator klastra
**Opis**: Wyjście funkcji powinno zawierać wyłącznie indexy spełniajace kryteria daty, źródła logów oraz identyfikatora klasta.

<figure>
    <img src="/reports/tests/reports-filter-indicies-test.png">
    <figcaption>Rysunek 1: Test filtrowania indexów ES na podstawie daty, źródła logów oraz identyfikatora klastra [źródło opracowanie własne]</figcaption>
</figure>

**Id**: TC2
**Tytuł**: Rozdzielanie logów na pakiety o określonym rozmiarze
**Opis**: Paczki logów są dzielone w taki sposób aby żadna z nich nie przekraczała określonego rozmiaru. Rozmiar paczki jest określany na podstawie długości zserializowanych logów w formacie JSON.

<figure>
    <img src="/reports/tests/reports-split-logs-into-packets-test.png">
    <figcaption>Rysunek 9: Testy podziału logów na pakiety [źródło opracowanie własne]</figcaption>
</figure>

**Id**: TC3
**Tytuł**: Serializacja tekstu do formatu **JSONL**
**Opis**: Otrzymując listę struktur, funkcja powinna zserializować je do formatu **JSON** oraz porozdzielać je znakiem nowej linii.

<figure>
    <img src="/reports/tests/reports-jsonl-encoder-test.png">
    <figcaption>Rysunek 1: Test kodera JSONL [źródło opracowanie własne]</figcaption>
</figure>

**Id**: TC4
**Tytuł**: Deserializacja **JSONL** do struktur w Go
**Opis**: Otrzymując tekst zawierający zseralizowane stuktury w formacie **JSONL**, funkcja powinna je zdeserializować i zwrócić listę struktur

<figure>
    <img src="/reports/tests/reports-jsonl-decoder-test.png">
    <figcaption> Test dekodera JSONL [źródło opracowanie własne]</figcaption>
</figure>

**Id**: TC5
**Tytuł**: Pakietowanie żądań do OpenAI API na podstawie ich wielkości
**Opis**: Żądanie zserializowane do formatu **JSONL** są dzielone na podstawie ich wielkości do pakietów nie przekraczających z góry określonego rozmiaru

<figure>
    <img src="/reports/tests/reports-client-split-test.png">
    <figcaption>Rysunek 1: Test podziału żądań do OpenAI API [źródło opracowanie własne]</figcaption>
</figure>

**Id**: TC6
**Tytuł**: Filtrowania logów aplikacji na podstawie dokładności
**Opis**: Filtr powinnien na wyjściu zwrócić jedynie logi, które zawierają słowa kluczowe odpowiadające poziomowi dokładności dla aplikacji, z której te logi zostały wyprodukowane.

<figure>
    <img src="/reports/tests/reports-application-filter-test.png">
    <figcaption> Test filtrowania logsów z aplikacji [źródło opracowanie własne]</figcaption>
</figure>

**Id**: TC7
**Tytuł**: Filtrowania logów aplikacji na podstawie dokładności
**Opis**: Filtr powinnien na wyjściu zwrócić jedynie logi, które zawierają słowa kluczowe odpowiadające poziomowi dokładności dla aplikacji, z której te logi zostały wyprodukowane.

<figure>
    <img src="/reports/tests/reports-node-filter-test.png">
    <figcaption>Rysunek 3: Test filtrowania logów z hostów [źródło opracowanie własne]</figcaption>
</figure>

**Id**: TC8
**Tytuł**: Wydobywanie **obserwacji** z odpowiedzi od OpenAI Batch API
**Opis**: Generator Obserwacji z hostów, deserializuje odpowiedź od OpenAI oraz dodaje do nich metadane aby utworzyć z nich **obserwacje** dla hosta

<figure>
    <img src="/reports/tests/reports-get-node-insights-from-batches-test.png">
    <figcaption>Rysunek 5: Test uzyskiwania obserwacji z hostów [źródło opracowanie własne]</figcaption>
</figure>

**Id**: TC9
**Tytuł**: Wydobywanie **obserwacji** z odpowiedzi od OpenAI Batch API
**Opis**: Generator Obserwacji z aplikacji, deserializuje odpowiedź od OpenAI oraz dodaje do nich metadane aby utworzyć z nich **obserwacje** dla aplikacji

<figure>
    <img src="/reports/tests/reports-get-application-insights-from-batches-test.png">
    <figcaption>Rysunek 6: Test uzyskiwania obserwacji z aplikacji [źródło opracowanie własne]</figcaption>
</figure>

**Id**: TC10
**Tytuł**: Tworznie raportów aplikacji na podstawie obserwacji dla aplikacji
**Opis**: Serwis Raportów tworzy raport aplikacji na podstawie obserwacji z aplikacji, grupując je na podstawie konkretnych aplikacji z jakich pochodzą oraz dodając metadane z jakimi dana obserwacja była wygenerowana.

<figure>
    <img src="/reports/tests/reports-application-reports-from-insights-test.png">
    <figcaption>Rysunek 7: Test tworzenia raportów aplikacji [źródło opracowanie własne]</figcaption>
</figure>

**Id**: TC11
**Tytuł**: Tworznie raportów hostów na podstawie obserwacji dla hostów
**Opis**: Serwis Raportów tworzy raport hostów na podstawie obserwacji z hostów, grupując je na podstawie konkretnych hostów z jakich pochodzą oraz dodając metadane z jakimi dana obserwacja była wygenerowana.

<figure>
    <img src="/reports/tests/reports-node-reports-from-insights-test.png">
    <figcaption>Rysunek 8: Test tworzenia raportów hostów [źródło opracowanie własne]</figcaption>
</figure>

**Id**: TC12
**Tytuł**: Uzyskiwanie pilności raportu
**Opis**: Pilność raportu jest ustalana na podstawie najwyższej pilności z wszystkich incydentów w raporcie

<figure>
    <img src="/reports/tests/reports-get-report-urgency-test.png">
    <figcaption>Rysunek 9: Test określania pilności raportu [źródło opracowanie własne]</figcaption>
</figure>

#### 9.1.2 Testy integracyjne {#testy-integracyjne-reports-service}

**Id**: TC13
**Tytuł**: Planowanie raportu
**Opis**: Planowanie raportu polega na zebraniu logów na podstawie źródeł sprecyzowanych w konfiguracji raportu oraz zakolejkowanie żądań do wygenerowania obserwacji z aplikacji i hostów.

**Kroki testowe**

1. Wypełnienie bazy logów aplikacji i hostów logami potrzebnymi do wygenerowania raportu
2. Przekazanie konfiguracji do funkcji odpowiedzialnej za planowanie raportu
3. Sprawdzanie stanu bazy raportów w poszukiwaniu wygenerowanego raportu

<figure>
    <img src="/reports/tests/reports-schedule-report-test.png">
    <figcaption>Rysunek 1: Test planowania raportu [źródło opracowanie własne]</figcaption>
</figure>

**Id**: TC14
**Tytuł**: Nasłuchiwanie na żądania wygenerowania raportu
**Opis**: Serwis nasłuchuje na żądania wygenerowania raportu z brokera wiadomości i po otrzymaniu takiej wiadomości tworzy zaplanowany raport

**Kroki testowe**

1. Wypełnienie bazy logów aplikacji i hostów logami potrzebnymi do wygenerowania raportu
2. Wysłanie wiadomości na brokera wiadomości odpowiedzialnego za przekazywanie żądań o wygenerowanie raportu
3. Oczekiwanie na otrzymanie i przetworzenie żądania przez serwis
4. Sprawdzanie stanu bazy raportów w poszukiwaniu wygenerowanego raportu

<figure>
    <img src="/reports/tests/reports-listen-for-report-requests.png">
    <figcaption> Test nasłuchiwania na żądania raportów [źródło opracowanie własne]</figcaption>
</figure>

#### 9.1.3 Pokrycie testów

Testami zostały pokryte kluczowe metody zawierające niebanalną logikę biznesową, metody te znajdowały się głównie w `/services`, `/incident_correlation`, `/insights` oraz `/handlers`

| Pakiet                                                                             | Pokrycie |
| ---------------------------------------------------------------------------------- | -------- |
| github.com/Magpie-Monitor/magpie-monitor/services/reports/internal/brokers         | 0.0%     |
| github.com/Magpie-Monitor/magpie-monitor/services/reports/cmd/reports              | 0.0%     |
| github.com/Magpie-Monitor/magpie-monitor/services/reports/internal/database        | 0.0%     |
| github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/config               | 0.0%     |
| github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/scheduled_jobs       | 0.0%     |
| github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/repositories         | 0.0%     |
| github.com/Magpie-Monitor/magpie-monitor/services/reports/internal/handlers        | 56.4%    |
| github.com/Magpie-Monitor/magpie-monitor/services/reports/internal/services        | 18.3%    |
| github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/filter               | 80.6%    |
| github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/incident_correlation | 30.2%    |
| github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/insights             | 20.7%    |
| github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/openai               | 5.8%     |

### 9.2 Testy Logs Ingestion Service {#testy-logs-ingestion-service}

#### 9.2.1 Testy integracyjne {#testy-integracyjne-logs-ingestion-service}

**Id**: TC1
**Tytuł**: Nasłuchiwanie na przychodzące z brokera logi hostów
**Opis**: Serwis nasłuchuje na przychodzące logi hostów i po otrzymaniu ich zapisuje je w odpowienim indeksie ElasticSearch

**Kroki testowe**

1. Wysłanie wiadomości z logami z hostów na brokera wiadomości
2. Oczekiwanie na odebranie wiadomości przez serwis
3. Sprawdzanie stanu bazy logów w poszukiwaniu nowych logów z hostów

<figure>
    <img src="/logs-ingestion/logs-ingestion-node-stream-reader-test.png">
    <figcaption>Rysunek 1: Test nasłuchiwania na wiadomości z logami hsotów [źródło opracowanie własne]</figcaption>
</figure>

**Id**: TC2
**Tytuł**: Nasłuchiwanie na przychodzące z brokera logi aplikacji
**Opis**: Serwis nasłuchuje na przychodzące logi aplikacji i po otrzymaniu ich transformuje i zapisuje je w odpowienim indeksie ElasticSearch

**Kroki testowe**

1. Wysłanie wiadomości z logami z aplikacji na brokera wiadomości
2. Oczekiwanie na odebranie wiadomości przez serwis
3. Sprawdzanie stanu bazy logów w poszukiwaniu nowych logów z aplikacji
<figure>
    <img src="/logs-ingestion/logs-ingestion-application-stream-reader-test.png">
    <figcaption> Test nasłuchiwania na wiadomości z logami aplikacji [źródło opracowanie własne]</figcaption>
</figure>

#### 9.2.1 Pokrycie testów

Testami były objęte głównie funkcje, które zawierają niebanalną logikę biznesową, która znajdowała się w pakiecie `logsstream`

| Pakiet                                                                              | Pokrycie |
| ----------------------------------------------------------------------------------- | -------- |
| github.com/Magpie-Monitor/magpie-monitor/services/logs_ingestion/pkg/config         | 0.0%     |
| github.com/Magpie-Monitor/magpie-monitor/services/logs_ingestion/cmd/logs_ingestion | 0.0%     |
| github.com/Magpie-Monitor/magpie-monitor/services/logs_ingestion/pkg/logs_stream    | 59.0%    |

## 9.3 Testy Agenta {#testy-agenta}

Agent przetestowany został jednostkowo w zakresie funkcjonalności rozdzielania zebranych logów na pakiety danych oraz deduplikacji logów, czyli procesu, w którym ze zbioru zebranych danych usuwane są logi, które powinny być częścią kolejnej paczki przesyłanych danych, a ich obecność w zbiorze wynika z niedokładności API klastra Kubernetes. W testach integracyjnych skupiono się natomiast na testowaniu integracji z API Kubernetes oraz zbieraniu logów z klastra.

### 9.3.1 Testy Node Agenta

#### 9.3.1.1 Testy jednostkowe

<figure>
    <img src="/agent/tests/agent-node-integration-example.png">
    <figcaption> Test odczytywania logów z pliku [źródło opracowanie własne]</figcaption>
</figure>

Przykładowy test przedstawia przyjęte podejście, permutacje danych wejściowych oraz wyjściowych są przekazywane do testu, który następnie weryfikuje poprawność działania funkcji.

**Id:** TC1  
**Tytuł:** Rozdzielanie zebranych logów na pakiety o podanym rozmiarze  
**Opis:**  
Agent rozdziela logi na pakiety danych o podanym rozmiarze.

**Warunki wstępne i oczekiwane rezultaty:**

<figure>
    <img src="/agent/tests/agent-node-unit-tc-1.png">
    <figcaption> Dane wejściowe oraz oczekiwane rezultaty testu podziału logów na pakiety [źródło opracowanie własne]</figcaption>
</figure>

Powyższe zdjęcie przedstawia wejściowe logi, rozmiar pakietu oraz spodziewaną liczbę pakietów dla wybranych przypadków testowych.

**Kroki testowe:**  
1\. Pobranie logów  
2\. Rozdzielenie logi na pakiety  
3\. Sprawdzenie czy liczba pakietów jest prawidłowa

#### 9.3.1.1 Testy integracyjne

<figure>
    <img src="/agent/tests/agent-node-integration-example.png">
    <figcaption> Przykładowy test obserwowania pliku z logami [źródło opracowanie własne]</figcaption>
</figure>

Przykładowy test przedstawia proces testowania czytania logów z pliku. W kontenerze tworzony jest plik, do którego następnie wpisywane są dane. Plik jest obserwowany przez Node Agenta, który odczytuje zapisane dane i przesyła rezultaty odczytu. Rezultaty są następnie porównywane z oczekiwanymi wynikami.

**Id:** TC2  
**Tytuł:** Odczytywanie logów z pliku przez Agenta  
**Opis:**  
Agent odczytuje logi z wybranego pliku.

**Warunki wstępne i oczekiwane rezultaty:**

<figure>
    <img src="/agent/tests/agent-node-integration-tc-1.png">
    <figcaption> Dane wejściowe oraz oczekiwane rezultaty testu obserwowania pliku z logami [źródło opracowanie własne]</figcaption>
</figure>

**Kroki testowe:**  
1\. Obserwowanie pliku przez agenta  
2\. Otworzenie pliku  
3\. Zapisanie przykładowych danych do pliku  
4\. Oczekiwanie na odczyt danych przez agenta  
5\. Porównanie odczytanych danych z oczekiwaniami

**Id:** TC3  
**Tytuł:** Zbieranie metadanych z hostów.  
**Opis:**  
Agent zbiera metadane o hoście na którym działa.

**Warunki wstępne i oczekiwane rezultaty:**

<figure>
    <img src="/agent/tests/agent-node-integration-tc-2.png">
    <figcaption> Dane wejściowe oraz oczekiwane rezultaty testu zbierającego metadane hostów [źródło opracowanie własne]</figcaption>
</figure>

**Kroki testowe:**  
1\. Utworzenie agenta z wejściową konfiguracją hosta  
2\. Odczytanie metadanych o hoście  
3\. Porównanie odczytanych wyników z oczekiwanymi rezultatami

### 9.3.2 Testy Pod Agenta

#### 9.3.2.1 Testy jednostkowe

<figure>
    <img src="/agent/tests/agent-pod-unit-example.png">
    <figcaption> Przykładowy test jednostkowy Pod Agenta [źródło opracowanie własne]</figcaption>
</figure>

Przykładowy test jednostkowy testuje rozdzielanie zebranych logów aplikacyjnych na pakiety kontenerów oraz podów wedle określonego rozmiaru. Logi są początkowo dzielone na pakiety kontenerów, które następnie grupowane są w pakiety podów o określonym rozmiarze.

**Id:** TC4  
**Tytuł:** Podział logów aplikacyjnych na pakiety kontenerów oraz podów  
**Opis:**  
Logi aplikacyjne przez agenta są dzielone na pakiety kontenerów oraz podów, wedle konfigurowalnego limitu rozmiaru.

**Warunki wstępne i oczekiwane rezultaty:**

<figure>
    <img src="/agent/tests/agent-pod-unit-tc-1.png">
    <figcaption> Dane wejściowe oraz oczekiwane rezultaty testu podziału logów na pakiety [źródło opracowanie własne]</figcaption>
</figure>

**Kroki testowe:**  
1\. Utworzenie agenta z wejściową konfiguracją rozmiaru pakietów  
2\. Podział logów na pakiety  
3\. Porównanie odczytanych wyników z oczekiwanymi rezultatami

**Id:** TC5  
**Tytuł:** Odczytanie wartości sekundy z timestampu RFC  
**Opis:**  
Agent powinien odczytywać wartość sekundy z timestampu RFC

**Warunki wstępne i oczekiwane rezultaty:**

<figure>
    <img src="/agent/tests/agent-pod-unit-tc-2.png">
    <figcaption> Dane wejściowe oraz oczekiwane rezultaty testu wyodrębnienia sekundy z timestampu RFC [źródło opracowanie własne]</figcaption>
</figure>

**Kroki testowe:**  
1\. Utworzenie agenta z wejściową konfiguracją rozmiaru pakietów  
2\. Podział logów na pakiety  
3\. Porównanie odczytanych wyników z oczekiwanymi rezultatami

**Id:** TC6  
**Tytuł:** Usuwanie logów z kolejnej sekundy  
**Opis:**  
Proces deduplikacji powinien usuwać logi z kolejnej sekundy, jeżeli takowe w zbiorze występują

**Warunki wstępne i oczekiwane rezultaty:**

<figure>
    <img src="/agent/tests/agent-pod-unit-tc-3.png">
    <figcaption> Dane wejściowe oraz oczekiwane rezultaty testu usuwania ze zbioru logów z kolejnej sekundy [źródło opracowanie własne]</figcaption>
</figure>

**Kroki testowe:**  
1\. Utworzenie agenta  
2\. Przeprowadzenie procesu usuwania logów  
3\. Porównanie odczytanych wyników z oczekiwanymi rezultatami

**Id:** TC7  
**Tytuł:** Podział logów z kontenerów na pakiety danych  
**Opis:**  
Logi z kontenerów powinny być dzielone na pakiety wedle zadanej konfiguracji

**Warunki wstępne i oczekiwane rezultaty:**

<figure>
    <img src="/agent/tests/agent-pod-unit-tc-4.png">
    <figcaption> Dane wejściowe oraz oczekiwane rezultaty testu podziału logów na pakiety [źródło opracowanie własne]</figcaption>
</figure>

**Kroki testowe:**  
1\. Utworzenie agenta  
2\. Przeprowadzenie procesu dzielenia logów na pakiety  
3\. Porównanie odczytanych wyników z oczekiwanymi rezultatami

#### 9.3.2.2 Testy integracyjne

W testach integracyjnych Pod Agenta skupiono się w głównej mierze na zbieraniu logów przy pomocy API Kubernetesa, które zamockowano. Mockowanie pozwoliło na całościowe przetestowanie procesu zbierania logów z klastra bez zewnętrznych zależności oraz zmian w kodzie.

<figure>
    <img src="/agent/tests/agent-pod-integration-example.png">
    <figcaption> Przykładowy test pobierania logów z zasobu StatefulSet [źródło opracowanie własne]</figcaption>
</figure>

<figure>
    <img src="/agent/tests/agent-pod-integration-example-1.png">
    <figcaption> Przykładowy test pobierania logów z zasobu StatefulSet - część 2 [źródło opracowanie własne]</figcaption>
</figure>

Przykładowy test integracyjny testuje zbieranie logów dla zasobu StatefulSet.

**Id:** TC8  
**Tytuł:** Zbieranie logów z zasobu Deployment  
**Opis:**  
Agent powinien zbierać logi zasobu Deployment zgodnie ze stanem klastra

**Warunki wstępne i oczekiwane rezultaty:**

<figure>
    <img src="/agent/tests/agent-pod-integration-tc-1.png">
    <figcaption> Dane wejściowe oraz oczekiwane rezultaty testu pobierania logów z zasobu Deployment [źródło opracowanie własne]</figcaption>
</figure>

**Kroki testowe:**  
1\. Utworzenie agenta z mockowym klientem Kubernetes oraz stanem wejściowym  
2\. Pobranie logów dla zasobu Deployment  
3\. Porównanie odczytanych wyników z oczekiwanymi rezultatami

**Id:** TC9  
**Tytuł:** Zbieranie logów z zasobu DaemonSet  
**Opis:**  
Agent powinien zbierać logi zasobu DaemonSet zgodnie ze stanem klastra

**Warunki wstępne i oczekiwane rezultaty:**

<figure>
    <img src="/agent/tests/agent-pod-integration-tc-2.png">
    <figcaption> Dane wejściowe oraz oczekiwane rezultaty testu pobierania logów z zasobu DaemonSet [źródło opracowanie własne]</figcaption>
</figure>

**Kroki testowe:**  
1\. Utworzenie agenta z mockowym klientem Kubernetes oraz stanem wejściowym  
2\. Pobranie logów dla zasobu DaemonSet  
3\. Porównanie odczytanych wyników z oczekiwanymi rezultatami

**Id:** TC10  
**Tytuł:** Zbieranie logów z zasobu StatefulSet  
**Opis:**  
Agent powinien zbierać logi zasobu StatefulSet zgodnie ze stanem klastra

**Warunki wstępne i oczekiwane rezultaty:**

<figure>
    <img src="/agent/tests/agent-pod-integration-tc-3.png">
    <figcaption> Dane wejściowe oraz oczekiwane rezultaty testu pobierania logów z zasobu StatefulSet [źródło opracowanie własne]</figcaption>
</figure>

**Kroki testowe:**  
1\. Utworzenie agenta z mockowym klientem Kubernetes oraz stanem wejściowym  
2\. Pobranie logów dla zasobu StatefulSet  
3\. Porównanie odczytanych wyników z oczekiwanymi rezultatami

**Id:** TC11  
**Tytuł:** Zbieranie danych o przestrzeniach nazw z klastra Kubernetes  
**Opis:**  
Agent powinien zbierać dane o przestrzeniach nazw z klastra Kubernetes

**Warunki wstępne i oczekiwane rezultaty:**

<figure>
    <img src="/agent/tests/agent-pod-integration-tc-4.png">
    <figcaption> Dane wejściowe oraz oczekiwane rezultaty testu pobierania danych o przestrzeniach nazw [źródło opracowanie własne]</figcaption>
</figure>

**Kroki testowe:**  
1\. Utworzenie agenta z mockowym klientem Kubernetes oraz stanem wejściowym  
2\. Pobranie danych o namespace  
3\. Porównanie odczytanych wyników z oczekiwanymi rezultatami

**Id:** TC12  
**Tytuł:** Zbieranie metadanych o aplikacjach z klastra Kubernetes  
**Opis:**  
Agent powinien zbierać dane o działających aplikacjach oraz ich rodzajach z klastra Kubernetes

**Warunki wstępne i oczekiwane rezultaty:**

<figure>
    <img src="/agent/tests/agent-pod-integration-tc-5.png">
    <figcaption> Dane wejściowe oraz oczekiwane rezultaty testu pobierania metadanych o aplikacjach [źródło opracowanie własne]</figcaption>
</figure>

**Kroki testowe:**  
1\. Utworzenie agenta z mockowym klientem Kubernetes oraz stanem wejściowym  
2\. Pobranie metadanych  
3\. Porównanie odczytanych wyników z oczekiwanymi rezultatami

#### 9.3.3 Pokrycie testów

| Pakiet                                                                                 | Pokrycie |
| -------------------------------------------------------------------------------------- | -------- |
| github.com/Magpie-Monitor/magpie-monitor/tree/main/agent/app/internal/agent/node/agent | 28.5%    |
| github.com/Magpie-Monitor/magpie-monitor/tree/main/agent/app/internal/agent/pods/agent | 48.3%    |

### 9.4 Testy Metadata Service {#testy-metadata-service}

Metadata Service został przetestowany integracyjnie, w obszarze pobierania metadanych z brokera Kafki oraz generowania zagregowanych metadanych.

#### 9.4.1 Testy integracyjne

**Id:** TC1  
**Tytuł:** Odbieranie metadanych o hostach z klastra Kubernetes  
**Opis:**  
Metadata service powinien odbierać metadane o hostach od Agenta

**Warunki wstępne i oczekiwane rezultaty:**

<figure>
    <img src="/metadata-service/tests/metadata-integration-tc-1.png">
    <figcaption> Dane wejściowe oraz oczekiwane rezultaty testu odbierania metadanych o hostach [źródło opracowanie własne]</figcaption>
</figure>

**Kroki testowe:**  
1\. Utworzenie instancji metadata service  
2\. Wysłanie przykładowych metadanych o hostach do brokera kafki  
3\. Porównanie odczytanych wyników z oczekiwanymi rezultatami

**Id:** TC2
**Tytuł:** Odbieranie metadanych o aplikacjach z klastra Kubernetes  
**Opis:**  
Metadata service powinien odbierać metadane o aplikacjach od Agenta

**Warunki wstępne i oczekiwane rezultaty:**

<figure>
    <img src="/metadata-service/tests/metadata-integration-tc-2.png">
    <figcaption> Dane wejściowe oraz oczekiwane rezultaty testu odbierania metadanych o aplikacjach [źródło opracowanie własne]</figcaption>
</figure>

**Kroki testowe:**  
1\. Utworzenie instancji metadata service  
2\. Wysłanie przykładowych metadanych o aplikacjach do brokera kafki  
3\. Porównanie odczytanych wyników z oczekiwanymi rezultatami

**Id:** TC3
**Tytuł:** Generowanie zagregowanych metadanych o hostach  
**Opis:**  
Metadata service powinien generować zagregowane metadane o hostach

**Warunki wstępne i oczekiwane rezultaty:**

<figure>
    <img src="/metadata-service/tests/metadata-integration-tc-3.png">
    <figcaption> Dane wejściowe oraz oczekiwane rezultaty testu generowania zagregowanych metadanych o hostach [źródło opracowanie własne]</figcaption>
</figure>

<figure>
    <img src="/metadata-service/tests/metadata-integration-tc-3-1.png">
    <figcaption> Dane wejściowe oraz oczekiwane rezultaty testu generowania zagregowanych metadanych o hostach - część 2 [źródło opracowanie własne]</figcaption>
</figure>

**Kroki testowe:**  
1\. Utworzenie instancji metadata service  
2\. Wysłanie przykładowych metadanych o hostach do brokera kafki  
3\. Nasłuchiwanie na wygenerowane zagregowane metadane o hostach  
4\. Porównanie odczytanych wyników z oczekiwanymi rezultatami

**Id:** TC4
**Tytuł:** Generowanie zagregowanych metadanych o aplikacjach  
**Opis:**  
Metadata service powinien generować zagregowane metadane o aplikacjach

**Warunki wstępne i oczekiwane rezultaty:**

<figure>
    <img src="/metadata-service/tests/metadata-integration-tc-4.png">
    <figcaption> Dane wejściowe oraz oczekiwane rezultaty testu generowania zagregowanych metadanych o aplikacjach [źródło opracowanie własne]</figcaption>
</figure>

<figure>
    <img src="/metadata-service/tests/metadata-integration-tc-4-1.png">
    <figcaption>Dane wejściowe oraz oczekiwane rezultaty testu generowania zagregowanych metadanych o aplikacjach - część 2 [źródło opracowanie własne]</figcaption>
</figure>

<figure>
    <img src="/metadata-service/tests/metadata-integration-tc-4-2.png">
    <figcaption>Dane wejściowe oraz oczekiwane rezultaty testu generowania zagregowanych metadanych o aplikacjach - część 3 [źródło opracowanie własne]</figcaption>
</figure>

**Kroki testowe:**  
1\. Utworzenie instancji metadata service  
2\. Wysłanie przykładowych metadanych o hostach do brokera kafki  
3\. Nasłuchiwanie na wygenerowane zagregowane metadane o aplikacjach  
4\. Porównanie odczytanych wyników z oczekiwanymi rezultatami

**Id:** TC5  
**Tytuł:** Generowanie zagregowanych metadanych o klastrach  
**Opis:**  
Metadata service powinien generować zagregowane metadane o klastrach

**Warunki wstępne i oczekiwane rezultaty:**

<figure>
    <img src="/metadata-service/tests/metadata-integration-tc-5.png">
    <figcaption> Dane wejściowe oraz oczekiwane rezultaty testu generowania zagregowanych metadanych o klastrach [źródło opracowanie własne]</figcaption>
</figure>

**Kroki testowe:**  
1\. Utworzenie instancji metadata service  
2\. Wysłanie przykładowych metadanych o hostach i aplikacjach do brokera kafki  
3\. Nasłuchiwanie na wygenerowane zagregowane metadane o klastrach  
4\. Porównanie odczytanych wyników z oczekiwanymi rezultatami

#### 9.4.2 Pokrycie testów

| Pakiet                                                                                                           | Pokrycie |
| ---------------------------------------------------------------------------------------------------------------- | -------- |
| github.com/Magpie-Monitor/magpie-monitor/blob/main/go/services/cluster_metadata/pkg/services/metadata_service.go | 76.8%    |

## 9.5. Testy Management Service {#testy-management-service}

Testy jednostkowe dla mikroserwisu `management-service` zostały zrealizowane przy użyciu języka Groovy oraz frameworka Spock. W celu oceny jakości testów i pokrycia kodu zastosowano narzędzie JaCoCo. Łączne pokrycie kodu wyniosło 54%, co obrazuje załączony wykres.

Testy zostały skoncentrowane na klasach zawierających logikę aplikacyjną, takich jak warstwy **service** i **utils**, które odgrywają kluczową rolę w przetwarzaniu danych. Klasy odpowiedzialne za zarządzanie kontrolerami, dostępem do danych, konfiguracją aplikacji oraz bezpieczeństwem nie zostały objęte testami jednostkowymi. Wynika to z faktu, że testowanie tych elementów wymagałoby podejścia integracyjnego/regresyjnego, wykraczającego poza zakres testów jednostkowych.

Podjęte działania pozwoliły na zweryfikowanie kluczowych funkcjonalności mikroserwisu.

<figure>
    <img src="/management-service/management-service-test-coverage.png">
    <figcaption> Pokrycie kodu testami [źródło opracowanie własne]</figcaption>
</figure>

| Pakiet                                                                               | Pokrycie |
| ------------------------------------------------------------------------------------ | -------- |
| github.com/Magpie-Monitor/magpie-monitor/management-service/src/main/java/pl/pwr/zpi | 54%      |

## 9.6 Testy funkcjonalne {#testy-funkcjonalne}

W celu weryfikacji czy system spełnia skonstruowane wymagania napisano scenariusze testowe:

### **Logowanie użytkownika**

**Opis**: Test ma na celu sprawdzenie, czy użytkownik może pomyślnie zalogować się do systemu, aby uzyskać dostęp do jego funkcjonalności.  
**Warunki wstępne**:

- Użytkownik posiada konto w systemie.

**Dane testowe**:

- Konto Google.

**Kroki testowe**:

- Użytkownik uruchamia aplikację Magpie Monitor.
- Użytkownik naciska przycisk “Sign in with Google”
- Użytkownik wybiera poprawne konto Google
- Użytkownik zostaje przekierowany na stronę główną systemu.

**Oczekiwany wynik**:

- Użytkownik pomyślnie loguje się do systemu i zostaje przeniesiony na stronę główną.
--------------------------------------------------------------------------
<br>

### **Próba dostępu przez nieuwierzytelnionego użytkownika**

**Opis**: Test ma na celu sprawdzenie, czy nieuwierzytelniony użytkownik może uzyskać dostęp do funkcjonalności systemu.  
**Warunki wstępne**:

- Użytkownik nie jest zalogowany w systemie.
- Użytkownik zna adres URL podstrony konfiguracji raportu

**Dane testowe**:

- Brak

**Kroki testowe**:

- Użytkownik wpisuje w przeglądarce adres URL strony z raportami

**Oczekiwany wynik**:

- Użytkownik nie otrzymuje dostępu do strony z raportami.
- Użytkownik zostaje przekierowany na stronę logowania

### **Wylogowanie użytkownika**

**Opis**: Test ma na celu sprawdzenie, czy użytkownik może poprawnie wylogować się z systemu.  
**Warunki wstępne**:

- Użytkownik jest zalogowany do systemu.

**Dane testowe**:

- Brak

**Kroki testowe**:

- Użytkownik jest zalogowany do systemu.
- Użytkownik klika przycisk "Wyloguj" w lewym dolnym rogu ekranu.
- Użytkownik zostaje przekierowany na stronę logowania.

**Oczekiwany wynik**:

- Użytkownik zostaje pomyślnie wylogowany i przekierowany na stronę logowania.

--------------------------------------------------------------------------
<br>

### **Wyświetlenie raportu**

**Opis**: Test ma na celu sprawdzenie, czy użytkownik może wyświetlić szczegóły raportu.  
**Warunki wstępne**:

- Użytkownik jest zalogowany do systemu.
- Istnieje wygenerowany raport z incydentami.
- Istnieje podłączony klaster do aplikacji

**Dane testowe**:

- Raport zawierający co najmniej jeden incydent.

**Kroki testowe**:

- Użytkownik loguje się do systemu.
- Użytkownik przechodzi do zakładki "Reports".
- Użytkownik klika na wybrany raport, aby zobaczyć szczegóły.
- Użytkownik przegląda stronę raportu.

**Oczekiwany wynik**:

- System wyświetla raport z nazwą klastra na podstawie, którego wygenerowano raport, przedział czasu, z którego zebrane zostały logi, statystyki, takie jak liczba przeanalizowanych aplikacji, hostów, liczba krytycznych, średnich oraz mało krytycznych incydentów, liczba przeanalizowanych logów z aplikacji, liczba przeanalizowanych logów z hostów.

--------------------------------------------------------------------------
<br>

### **Przeglądanie incydentów**

**Opis**: Test ma na celu sprawdzenie, czy użytkownik może wyświetlić listę incydentów wykrytych w wygenerowanym raporcie.  
**Warunki wstępne**:

- Użytkownik jest zalogowany do systemu.
- Istnieje wygenerowany raport z incydentami.
- Istnieje podłączony do systemu klaster z aplikacją.

**Dane testowe**:

- Raport zawierający co najmniej jeden incydent.

**Kroki testowe**:

- Użytkownik loguje się do systemu.
- Użytkownik przechodzi do zakładki "Reports".
- Użytkownik klika na wybrany raport, aby zobaczyć szczegóły.
- Użytkownik przegląda listę incydentów w raporcie.

**Oczekiwany wynik**:

- System wyświetla listę incydentów, zawierającą nazwę aplikacji/hosta, kategorię i tytuł incydentu, oraz datę wykrycia.

--------------------------------------------------------------------------
<br>

### **Generacja raportu na żądanie**

**Opis**: Test ma na celu sprawdzenie, czy użytkownik może wyświetlić listę incydentów wykrytych w wygenerowanym raporcie.  
**Warunki wstępne**:

- Użytkownik jest zalogowany do systemu.
- Istnieje wygenerowany raport z incydentami.
- Istnieje podłączony do systemu klaster z aplikacją.

**Dane testowe**:

- Raport zawierający co najmniej jeden incydent.

**Kroki testowe**:

- Użytkownik loguje się do systemu.
- Użytkownik przechodzi do zakładki "Reports".
- Użytkownik klika na wybrany raport, aby zobaczyć szczegóły.
- Użytkownik przegląda listę incydentów w raporcie.

**Oczekiwany wynik**:

- System wyświetla listę incydentów, zawierającą nazwę aplikacji/hosta, kategorię i tytuł incydentu, oraz datę wykrycia.

--------------------------------------------------------------------------
<br>

### **Konfiguracja dokładności aplikacji i hostów**

**Opis**: Test ma na celu sprawdzenie, czy użytkownik może dostosować poziomy dokładności osobno dla każdej aplikacji i hosta w raporcie.  
**Warunki wstępne**:

- Użytkownik jest zalogowany do systemu.
- Istnieje klaster podłączony do systemu.
- Na istniejącym klastrze występują przynajmniej dwa aplikacje
- Na istniejącym klastrze występują przynajmniej dwa hosty

**Dane testowe**:

- Klaster, zawierający co najmniej dwie aplikacje oraz dwóch hostów.

**Kroki testowe**:

- Użytkownik loguje się do systemu.
- Użytkownik przechodzi do zakładki "Clusters".
- Użytkownik naciska przycisk “Report configuration”
- Użytkownik dodaje dwie aplikacje
- Użytkownik zmienia dokładność pierwszej z nich na “low”
- Użytkownik dodaje dwóch hostów
- Użytkownik zmienia dokładność pierwszego z nich na “high”, a drugiego na “medium”

**Oczekiwany wynik**:

- System wyświetla listę dodanych aplikacji oraz hostów do raportu z różnymi wartościami pola dokładność

--------------------------------------------------------------------------
<br>

### **Planowanie raportów**

**Opis**: Test ma na celu sprawdzenie, czy użytkownik może zaplanować cykliczne generowanie raportów w określonych przedziałach czasu.  
**Warunki wstępne**:

- Użytkownik jest zalogowany do systemu.
- Istnieje klaster podłączony do systemu.

**Dane testowe**:

- Brak

**Kroki testowe**:

- Użytkownik loguje się do systemu.
- Użytkownik przechodzi do zakładki "Clusters".
- Użytkownik naciska przycisk “Report configuration”
- Użytkownik ustawia wartość ”Generation type” na “Scheduled”
- Użytkownik ustawia “Schedule period” na “1 week”
- Użytkownik naciska przycisk “Generate”
- Użytkownik odczekuje tydzień

**Oczekiwany wynik**:

- System wygenerował raport z ostatniego tygodnia
- W zakładce “Reports” pojawił się nowy zaplanowany raport na przyszły tydzień

--------------------------------------------------------------------------
<br>

### **Generacja raportu na żądanie**

**Opis**: Test ma na celu sprawdzenie, czy użytkownik może generować raport na żądanie z wybranego przedziału czasu.  
**Warunki wstępne**:

- Użytkownik jest zalogowany do systemu.
- Istnieje klaster podłączony do systemu.

**Dane testowe**:

- Data początkowa przedziału czasu, z którego ma powstać raport
- Data końcowa przedziału czasu, z którego ma powstać raport

**Kroki testowe**:

- Użytkownik loguje się do systemu.
- Użytkownik przechodzi do zakładki "Clusters".
- Użytkownik naciska przycisk “Report configuration”
- Użytkownik ustawia wartość ”Generation type” na “ON_DEMAND”
- Użytkownik ustawia “Data Range” na daty z danych testowych
- Użytkownik naciska przycisk “Generate”
- Użytkownik odczekuje aż raport zostanie przygotowany

**Oczekiwany wynik**:

- System wygenerował raport z wybranego przedziału
- Utworzony raport nie zawiera incydentu spoza wybranego zakresu czasu

--------------------------------------------------------------------------
<br>

### **Personalizacja interpretacji logów**

**Opis**: Test ma na celu sprawdzenie, czy użytkownik może dodać własne instrukcje dla modelu językowego, aby raporty były dostosowane do specyfiki aplikacji i hostów.  
**Warunki wstępne**:

- Użytkownik jest zalogowany do systemu.
- Istnieje klaster podłączony do systemu, zawierający minimum jedną aplikację i jednego hosta.

**Dane testowe**:

- Instrukcja dostosowująca interpretację logów w raporcie.

**Kroki testowe**:

- Użytkownik loguje się do systemu.
- Użytkownik przechodzi do zakładki "Clusters".
- Użytkownik naciska przycisk “Report configuration”
- Użytkownik ustawia wartość ”Generation type” na “ON_DEMAND”
- Użytkownik dodaje aplikację do raportu
- Użytkownik dodaje hosta do raportu
- Użytkownik dodaje własną interpretacje logów do aplikacji
- Użytkownik naciska przycisk “Generate”
- Użytkownik odczekuje aż raport zostanie przygotowany

**Oczekiwany wynik**:

- System wygenerował raport
- Incydenty aplikacji zawierają informacje o wprowadzonej instrukcji do interpretacji logów

--------------------------------------------------------------------------
<br>

### **Wybór analizowanych hostów**

**Opis**: Test ma na celu sprawdzenie, czy użytkownik może wskazać konkretne hosty, które będą źródłem danych do raportu.  
**Warunki wstępne**:

- Użytkownik jest zalogowany do systemu.
- Istnieje klaster podłączony do systemu, zawierający minimum dwa hosty.

**Dane testowe**:

- Lista dostępnych hostów dla istniejącego klastra

**Kroki testowe**:

- Użytkownik loguje się do systemu.
- Użytkownik przechodzi do zakładki "Clusters".
- Użytkownik naciska przycisk “Report configuration”
- Użytkownik ustawia wartość ”Generation type” na “ON_DEMAND”
- Użytkownik dodaje pierwszego hosta do raportu
- Użytkownik naciska przycisk “Generate”
- Użytkownik odczekuje aż raport zostanie przygotowany

**Oczekiwany wynik**:

- System wygenerował raport
- Każdy znaleziony incydent dotyczy tylko wybranego hosta

--------------------------------------------------------------------------
<br>

### **Wybór analizowanych aplikacji**

**Opis**: Test ma na celu sprawdzenie, czy użytkownik może wskazać konkretne aplikacje, które będą źródłem danych do raportu.  
**Warunki wstępne**:

- Użytkownik jest zalogowany do systemu.
- Istnieje klaster podłączony do systemu, zawierający minimum dwie aplikacje.

**Dane testowe**:

- Lista dostępnych aplikacji dla istniejącego klastra

**Kroki testowe**:

- Użytkownik loguje się do systemu.
- Użytkownik przechodzi do zakładki "Clusters".
- Użytkownik naciska przycisk “Report configuration”
- Użytkownik ustawia wartość ”Generation type” na “ON_DEMAND”
- Użytkownik dodaje pierwszą aplikację do raportu

- Użytkownik naciska przycisk “Generate”
- Użytkownik odczekuje aż raport zostanie przygotowany

**Oczekiwany wynik**:

- System wygenerował raport
- Każdy znaleziony incydent dotyczy tylko wybranej aplikacji

--------------------------------------------------------------------------
<br>

### **Dodanie kanału powiadomień do raportu**

**Opis**: Test ma na celu sprawdzenie, czy użytkownik może przypisać kanał powiadomień (Slack, Discord lub email) do raportu.  
**Warunki wstępne**:

- Użytkownik jest zalogowany do systemu.
- Istnieje klaster podłączony do systemu, zawierający minimum dwie aplikacje.

**Dane testowe**:

- Dostępny kanał powiadomień

**Kroki testowe**:

- Użytkownik loguje się do systemu.
- Użytkownik przechodzi do zakładki "Clusters".
- Użytkownik naciska przycisk “Report configuration”
- Użytkownik ustawia wartość ”Generation type” na “ON_DEMAND”
- Użytkownik dodaje do dostępny kanał powiadomień do raportu
- Użytkownik naciska przycisk “Generate”
- Użytkownik odczekuje aż raport zostanie przygotowany

**Oczekiwany wynik**:

- System wygenerował raport
- W momencie wygenerowania raportu, wiadomość o zakończeniu pracy przychodzi na przypisany kanał powiadomień

--------------------------------------------------------------------------
<br>

### **Dodanie nowego kanału powiadomień**

**Opis**: Test ma na celu sprawdzenie, czy użytkownik dodać nowy kanał powiadomień (Slack, Discord lub email) do raportu.  
**Warunki wstępne**:

- Użytkownik jest zalogowany do systemu.

**Dane testowe**:

- Dostępny kanał powiadomień

**Kroki testowe**:

- Użytkownik loguje się do systemu.
- Użytkownik przechodzi do zakładki "Notifications".
- Użytkownik dodaje kanał powiadomień do odpowiedniej kategorii

**Oczekiwany wynik**:

- Do systemu został dodany nowy kanał powiadomień
- Kanał powiadomień zostaje wyświetlany na stronie “Notifications” z poprawną nazwą, webhookiem oraz datami dodania oraz modyfikacji

--------------------------------------------------------------------------
<br>

### **Usuwanie kanału powiadomień**

**Opis**: Test ma na celu sprawdzenie, czy użytkownik może usunąć kanał powiadomień.  
**Warunki wstępne**:

- Użytkownik jest zalogowany do systemu.

**Dane testowe**:

- Dodany wcześniej kanał powiadomień, który zostanie usunięty

**Kroki testowe**:

- Użytkownik loguje się do systemu.
- Użytkownik przechodzi do zakładki "Notifications".
- Użytkownik wybiera kanał powiadomień do usunięcia.
- Użytkownik wybiera ikonkę “Kosz”, aby usunąć kanał

**Oczekiwany wynik**:

- Z systemu zostaje usunięty kanał powiadomień

--------------------------------------------------------------------------
<br>

### **Edytowanie kanału powiadomień**

**Opis:** Test ma na celu sprawdzenie, czy użytkownik może edytować dodany kanał powiadomień.  
**Warunki wstępne**:

- Użytkownik jest zalogowany do systemu.

**Dane testowe**:

- Dodany wcześniej kanał powiadomień (np. adres email lub kanał Slack) do edycji.
- Nowa nazwa kanału powiadomień

**Kroki testowe**:

- Użytkownik loguje się do systemu.
- Użytkownik przechodzi do zakładki "Notifications".
- Użytkownik wybiera kanał powiadomień do edycji.
- Użytkownik wybiera ikonkę “Ołówka”, aby usunąć kanał.
- W nowo otwartym okienku użytkownik zmienia nazwę kanału.
- Użytkownik naciska przycisk “Save”.

**Oczekiwany wynik**:

- Kanał powiadomień zmienia nazwę na nową.

--------------------------------------------------------------------------
<br>

### **Testowanie powiadomień**

**Opis:** Test ma na celu sprawdzenie, czy użytkownik może przetestować dodany kanał powiadomień.  
**Warunki wstępne**:

- Użytkownik jest zalogowany do systemu.
- Użytkownik ma uprawnienia do wyświetlania wiadomości na wybranym kanale powiadomień.

**Dane testowe**:

- Dodany wcześniej kanał powiadomień (np. adres email lub kanał Slack).

**Kroki testowe**:

- Użytkownik loguje się do systemu.
- Użytkownik przechodzi do zakładki "Notifications".
- Użytkownik wybiera kanał powiadomień do przetestowania.
- Użytkownik wybiera ikonę “Wysłania wiadomości”.

**Oczekiwany wynik**:

- Na kanale powiadomień pojawia się nowa wiadomość testowa wysłana przez Magpie Monitor.

--------------------------------------------------------------------------
<br>

### **Wyświetlenie danych o incydencie**

**Opis:** Test ma na celu sprawdzenie, czy użytkownik może wyświetlić dane o incydencie.  
**Warunki wstępne**:

- Użytkownik jest zalogowany do systemu.
- Istnieje wygenerowany raport z co najmniej jednym incydentem.

**Dane testowe**:

- Raport, zawierający co najmniej jeden incydent.

**Kroki testowe**:

- Użytkownik loguje się do systemu.
- Użytkownik przechodzi do zakładki "Reports".
- Użytkownik wybiera raport z incydentem.
- Użytkownik przechodzi na stronę incydentu.
- Użytkownik znajduję wszystkie dane dotyczące incydentu.

**Oczekiwany wynik**:

- Użytkownik powinien móc wyświetlić: nazwę aplikacji lub hosta; rekomendację działań opis; czas; w jakim doszło do incydentu.

--------------------------------------------------------------------------
<br>

### **Wyświetlenie listy raportów**

**Opis:** Test ma na celu sprawdzenie, czy użytkownik może wyświetlić listę raportów.  
**Warunki wstępne**:

- Użytkownik jest zalogowany do systemu.
- Istnieje wygenerowany raport na żądanie.
- Istnieje wygenerowany raport cykliczny.
- Istnieją zaplanowane raporty do generacji.

**Dane testowe**:

- Wygenerowany raport na żądanie.
- Wygenerowany raport cykliczny.
- Raport, który jest w trakcie generacji.

**Kroki testowe**:

- Użytkownik loguje się do systemu.
- Użytkownik przechodzi do zakładki "Reports".

**Oczekiwany wynik**:

- Użytkownik powinien móc zobaczyć listę raportów, z podziałem na raporty cykliczne oraz na żądanie.
- Wśród raportów powinien widnieć generowany raport wraz z informacją, że proces jego generowania jest w toku.

--------------------------------------------------------------------------
<br>

# 10. Podsumowanie {#podsumowanie}

## 10.1 Przebieg projektu {#przebieg-projektu}

Projekt był realizowany w metodyce zwinnej. Na początkowym etapie szczególna uwaga była przyłożona do projektu architektury oraz interfejsu użytkownika, tak aby zdecydować się na odpowiednie rozwiązania technologiczne oraz przygotować interfejsy na wymagania funkcjonalne, które system ma oferować.

Po realizacji wstępnej dokumentacji i doboru technologi, rozpoczęta została praca implementacyjna - **2 sprint**. Zadania zostały zorganizowamy w taki sposób aby móc jak najszybciej rozpocząć równoległą pracę nad wszystkimi mikroserwisami w systemie.

Na początkowym etapie projektu, kluczowa była również weryfikacja oczekiwać odnośnie analizy dużej ilości logów przez model językowy. W związku z tym w ramach pierwszych sprintów, kluczowe było zaimplementowanie podstawowych funkcjonalności Logs Ingestion Service oraz Reports Service, ponieważ wykrycie problemów z tym podsystemem mogło uniemożliwić realizację całego projektu. Na wczesnym etapie zweryfikowane zatem było, że podstawowe założenia projektu mają sens zarówno techniczny jak i biznesowy.

W ramach **sprintu 4** oraz **sprintu 5**, udało się zrealizować integracje mikroserwisów oraz wykonać testy obciążeniowe, które zakończyły się sukcesem.

Ostatnimi zadaniami w projekcie były naprawa drobnych błędów oraz rozpisanie testów jednostkowych i integracyjnych, co zostało wykonane w **spincie 7** oraz **sprincie 8**.

### 10.1.1 Ós czasu

Projekt był realizowany w ramach 8 sprintów. Pierwsze 2 sprinty trwały tydzień ze względu na dużą dynamikę oraz zmieniające się wymagania i koncepty na początku życia projektu. Następne sprinty trwały standardowe 2 tygodnie. Po każdym sprincie, wszyscy członkowie zespołu planowali zadania na następny sprint i przygotowywali estymaty punktów story point, tak aby nowe zadania były możliwe do realizacji w trakcie następnego sprintu.

<figure>
    <img src="/przebieg-projektu/jira-os-czasu.png">
    <figcaption> Oś czasu w wykonywaniu projektu [źródło opracowanie własne]</figcaption>
</figure>

### 10.1.2 Diagram zaangażowania

Diagram zaangażowania pozwala zobaczyć ile zadań udało się wykonać w ramach danego sprintu w kontekście wszystkich przypisanych zadań.
W większości sprintów pojedyńcze zadania musiały być przenoszone do następnych sprintów. Te decyzje były podjęte ze względu na pojawiające się w trakcie sprintu zadania, które musiały być jak najszybciej wykonane.

<figure>
    <img src="/przebieg-projektu/jira-zaangazowanie.png">
    <figcaption> Wykres zaangażowania [źródło opracowanie własne]</figcaption>
</figure>

### 10.1.3 Diagramy spalania

Diagramy spalania wskazują ile punktów story point było wykonywanych w trakcie trwania sprintu. Punktem odniesienie jest przewodnik, czyli sugestia odnośnie idealnego tempa prac w projekcie. Niestety diagram ten nie uwzględnia, że wiele zadań było dopisywanych w trakcie sprintu. Zadania te były dopisywane razem z nowymi wymagania, które musiałby być pilnie zrealizowane.

<figure>
    <img src="/przebieg-projektu/sprint-3-burndown.png">
    <figcaption> Wykres spalania ze Sprintu 3. [źródło opracowanie własne]</figcaption>
</figure>

<figure>
    <img src="/przebieg-projektu/sprint-4-burndown.png">
    <figcaption> Wykres spalania ze Sprintu 4. [źródło opracowanie własne]</figcaption>
</figure>

<figure>
    <img src="/przebieg-projektu/sprint-5-burndown.png">
    <figcaption> Wykres spalania ze Sprintu 5. [źródło opracowanie własne]</figcaption>
</figure>

<figure>
    <img src="/przebieg-projektu/sprint-6-burndown.png">
    <figcaption> Wykres spalania ze Sprintu 6. [źródło opracowanie własne]</figcaption>
</figure>

<figure>
    <img src="/przebieg-projektu/sprint-7-burndown.png">
    <figcaption> Wykres spalania ze Sprintu 7. [źródło opracowanie własne]</figcaption>
</figure>

<figure>
    <img src="/przebieg-projektu/sprint-8-burndown.png">
    <figcaption> Wykres spalania ze Sprintu 8. [źródło opracowanie własne]</figcaption>
</figure>

### 10.1.4 Statystyki systemu kontroli wersji

System kontroli wersji wykorzystany w projekcie - git wraz z dostawcą, który oferuje zdalną synchronizację - Github, oferuje statystyki, które pozwalają na bardziej realistyczną estymate postępów prac.

#### 10.1.4.1 Liczba kontrybucji

W momencie tworzenia dokumentacji w projekcie wykonano **485 kontrybucji**. Każda z nich podegała wymogom odnośnie opisów i treści. Efektem takiego podejścia było powstanie kontroli wersji w której łatwo można analizować i rozumieć zmiany w kodzie za pośrednictwem wiadomości w kontrybucjach. Dodatkowo operacje takie jak wybiórcze dodawawanie zmian (cherry-pick) lub wycofywanie zmian (revert) stały się również bardzo przystępne.

<figure>
    <img src="/przebieg-projektu/github-commits-number.png">
    <figcaption> Liczba kontrybucji w projekcie [źródło opracowanie własne]</figcaption>
</figure>

Kontrybucje mają zawsze postać: **{nazwa-częsci-projektu}: Wykonana zmiana**

<figure>
    <img src="/przebieg-projektu/github-commits-format.png">
    <figcaption> Format kontrybucji [źródło opracowanie własne]</figcaption>
</figure>

W trakcie realizacji projektu zostały zamknięte **137 pull requesty**. Pull request odpowiada w większości przypadków dodaniu nowych funkcjonalności lub naprawienia obecnych błedów. Każdy pull request wymagał akceptacji innego członka zespołu. W więkości przypadków recenzenci dodawali komentarze sugerujące poprawki i dopiero po ich zaimplementowaniu akceptowali zmiany.

<figure>
    <img src="/przebieg-projektu/github-prs.png">
    <figcaption> Liczba Pull Requestów [źródło opracowanie własne]</figcaption>
</figure>

#### 10.1.4.2 Przebieg projektu według liczby kontrybucji

Na poniższym wykresie można zobaczyć liczbę liń kodu dodano do głównej gałęzi w trakcie realizacji projektu.

<figure>
    <img src="/przebieg-projektu/additions-graph.png">
    <figcaption> Wykres dodanych linii kodu w projekcie [źródło opracowanie własne]</figcaption>
</figure>

Poniższy wykres przedstawia liczbę kontrybucji (commitów) do głównej gałęzi repozytorium

<figure>
    <img src="/przebieg-projektu/commits-graph.png">
    <figcaption>  Wykres kontrybucji w projekcie [źródło opracowanie własne]</figcaption>
</figure>

## 10.2 Wnioski {#wnioski}

### 10.2.1 Wnioski z konceptu projektu

Projekt został zaplanowany jako dowód konceptu (Proof of concept) zintegrowanego systemu do zbierania i przetwarzania logów z dużych klastów komputerowych, i w ramach tego celu udało się zweryfikować, że rozwiązanie takie jest nie tylko możliwe technicznie, lecz również posiada sens biznesowy.

### 10.2.2 Wnioski z implementacji projektu

Podczas pracy z modelem oferowanym przez OpenAI oraz pracy z Batch API od tego samego dostawcy odkryto, że rozwiązanie to pomimo atrakcyjnego modelu biznesowego (dostosowanego do dużej ilości danych), nie jest rowiązaniem bezproblemowych.

W trakcie implementacji odkryto nieprzewidywalne odkładanie zadań przez Batch API co powodowało utrudnione testowanie nowych funkcjonalności end-2-end. Dodatkowo, na niektóre żądania Batch API od OpenAI zwracał błąd z kodem **500** bez podania przyczyny.
Taki kod błędu wskazuje na wewnętrzny błąd systemu i może sugerować niestabilność oferowanego przez OpenAI produktu.

Zaobserwowano również jak dużą kontrole nad dokładnością obserwacji modelu ma rozmiar danych umieszczonych w jednym kontekście. Zmniejszając tą wartość możliwe było uzyskiwanie lepszych obserwacji kosztem większego kosztu.

Jednym z głównych potencjałów na wartość intelektualną w tym rozwiązaniu było filtrowanie logów przed przetwarzaniem przez model językowy. Perfekcyjne rozwiązanie minimalizuje koszty przetwarzania logów, jednocześnie nie zmniejszając ich dokładności.

Wypracowane rozwiązanie, opierające się na bogatym, domenowych zbiorze słów kluczowych pozwaliło na drastyczną redukcje kosztów. W tym przypadku nietypowym paramerem, który dodatkowo pozwalał na kontrolę zachowania filtra były maksymalne rozmiary pojedyńczych logów. Im wyższy maksymalny rozmiar loga tym filtr bardziej ograniczał ostateczną liczbę logów.

Filtr ma potencjał na wykorzystanie tańszego modelu do analizy istotności loga i filtrowanie go przy jego użyciu. Potencjalnymi rozwiązaniami tego problemu są również klasyfikatory lub analizatory sentymentu.

Roszerzeniem projektu mający duży potencjał na obiniżenie kosztów i dodatkowe usprawnienia w obserwacji logów oferuje własny model językowy, który byłby mniejszy ale lepiej dotrenowany do konkretnego zadania, takiego jak wykrywanie anomalii, tworzenie podsumowań czy rekomendację rozwiązań.

Mniej domenowym aspektem, który był kluczowy podczas realizacji Magpie Monitor, było wykorzystywanie architektury mikroserwisowej zgodnie z konceptem event-driven (zorientowanej na wydarzenia). Pozwoliło to na drastyczne zwiększenie odporności na awarie przez usunięcie powiązań pomiędzy mikroserwisami oraz zwiększenie potencjału na skalowanie zarówno horyzontalne jak i wertykalne.

# 11. Dokumentacja użytkownika {#dokumentacja-użytkownika}

## 11.1 Wprowadzenie {#wprowadzenie}

W celu ułatwienia użytkownikom rozpoczęcia pracy z Magpie Monitorem, przygotowano instrukcję opisującą, jak uzyskać dostęp do podstawowych funkcji systemu oraz skutecznie z nich korzystać.  
Magpie Monitor to zaawansowane narzędzie do monitorowania logów pochodzących z wybranego klastra Kubernetesa. W związku z tym instalacja wymaga dodania specjalnego serwisu do istniejącej architektury. Szczegółowy opis procesu instalacji znajduje się w kolejnej sekcji [instalacja aplikacji](#instalacja-aplikacji).  
Ponadto dokumentacja obejmuje omówienie kluczowych scenariuszy użytkowania, takich jak:

- [logowanie do aplikacji](#logowanie-do-aplikacji),
- [otworzenie ostatniego raportu](#otworzenie-ostatniego-raportu),
- [planowanie generowania raportów](#planowanie-generowania-raportów),
- [generacja raportu na życzenie](#generacja-raportu-na-życzenie),
- [konfiguracja kanałów powiadomień](#konfiguracja-kanałów-powiadomień)

Przedstawiona instrukcja stanowi solidną podstawę do zapoznania się z najważniejszymi funkcjami systemu, umożliwiając efektywne wykorzystanie jego możliwości w codziennym monitorowaniu logów.


## 11.2 Użytkowanie produktu programowego {#użytkowanie-produktu-programowego}


### 11.2.1 Instalacja aplikacji {#instalacja-aplikacji}

Aby zacząć korzystać z systemu Magpie Monitor, należy zainstalować na swoim klastrze komputerowym Kubernetes i odpowiednio skonfigurować agenta, który zbiera logi z aplikacji oraz hostów.

<figure>
    <img src="/agent/agent-helm-values.png">
    <figcaption> Konfiguracja paczki wdrożeniowej Helm Agenta [źródło opracowanie własne]</figcaption>
</figure>

Agent dostarczany jest w paczce wdrożeniowej Helm, której konfiguracja zawiera się w pliku values.yaml. Klient może skonfigurować agenta wedle swoich potrzeb, dostępne opcje konfiguracyjne to m.in.:

- przyjazna nazwa klastra, będąca jego identyfikatorem
- wyłączone ze zbierania logów przestrzenie nazw klastra
- pliki, z których zbierane są logi

Skonfigurowaną paczkę wdrożeniową instaluje sie przy pomocy komendy _helm install_, zgodnie z dokumentacją narzędzia Helm (https://helm.sh/docs/helm/helm_install/).

   <figure>
    <img src="/agent/agent-installed-view.png">
    <figcaption> Kontenery wdrożonego Agenta na klastrze Kubernetes [źródło opracowanie własne]</figcaption>
</figure>

Po zainstalowaniu, agent zacznie wysyłać logi oraz metadane do chmury Magpie Monitor, która rozpocznie analizę logów.

### 11.2.2 Najczęściej wykonywane operacje {#najczęściej-wykonywane-operacje}


#### 11.2.2.1 Logowanie do aplikacji {#logowanie-do-aplikacji}

Wymaganie funkcjonalne: REQ01

W celu zalogowania się do systemu należy otworzyć stronę Magpie Monitor. Następnie wybrać opcję “Sign in with Google”:

<figure>
    <img src="/user-interface/login-page.png">
    <figcaption> Widok logowania [źródło opracowanie własne]</figcaption>
</figure>

Po naciśnięciu przycisku użytkownik zostaje przeniesiony na stronę dostarczoną przez firmę Google:

<figure>
    <img src="/user-interface/google-sign-in-page.png">
    <figcaption> Widok logowania dostarczony przez Google  [źródło opracowanie własne]</figcaption>
</figure>

Po wybraniu odpowiedniego konta zostaniemy zalogowani i przeniesieni na widok główny.

#### 11.2.2.2 Otworzenie ostatniego raportu {#otworzenie-ostatniego-raportu}

Po zalogowaniu do systemu widok główny zawsze zawiera ostatni wygenerowany raport:

<figure>
    <img src="/user-interface/main-page.png">
    <figcaption> Widok główny [źródło opracowanie własne]</figcaption>
</figure>Przedstawiona instrukcja stanowi solidną podstawę do zapoznania się z najważniejszymi funkcjami Magpie Monitora, umożliwiając efektywne wykorzystanie jego możliwości w codziennym monitorowaniu logów.

Alternatywnie możemy również otworzyć ostatni raport z widoku raportów. Po kliknięciu zakładki “Reports”, użytkownik zostanie przeniesiony do tego widoku:

<figure>
    <img src="/user-interface/reports-page.png">
    <figcaption> Widok raportów [źródło opracowanie własne]</figcaption>
</figure>

Raporty są posortowane malejąco według daty rozpoczęcia procesu generacji. Na powyższym obrazku najnowszy raport pochodzi z dnia 06.12.2024. Po wyborze przycisku w sekcji “Actions” użytkownik zostaje przeniesiony do strony raportu:

<figure>
    <img src="/user-interface/report-page.png">
    <figcaption> Widok raportu [źródło opracowanie własne]</figcaption>
</figure>

### 11.2.2.3 Planowanie generowania raportów {#planowanie-generowania-raportów}

Wymagania funkcjonalne: REQ05, REQ06, REQ07, REQ09, REQ10, REQ11, REQ12   

Po zalogowaniu należy przejść do zakładki “Clusters”:

<figure>
    <img src="/user-interface/clusters-page.png">
    <figcaption> Widok klastrów [źródło opracowanie własne]</figcaption>
</figure>

Po wybraniu klastra, dla którego chcemy wygenerować raport cykliczny, należy nacisnąć przycisk “Report configuration”. Użytkownik zostanie przeniesiony wtedy do widoku konfiguracji raportu:

<figure>
    <img src="/user-interface/report-config-scheduled.png">
    <figcaption> Widok konfiguracji raportu dla raportów cyklicznych [źródło opracowanie własne]</figcaption>
</figure>

Aby wygenerować raport cykliczny, w sekcji „Generation type” należy wybrać opcję „Scheduled”. Po jej zaznaczeniu pojawi się sekcja „Schedule period”, w której można określić, w jakich odstępach czasu raport ma być generowany. Następnie należy skonfigurować pozostałe elementy raportu, takie jak kanały powiadomień, które mają informować o zakończeniu procesu generacji, oraz wybrać aplikacje i hosty, których logi mają zostać uwzględnione. Dla każdej aplikacji i hosta można również precyzyjnie dostosować instrukcje do modelu, klikając ikonę w kolumnie „Custom prompt”, oraz zmienić poziom dokładności w kolumnie „Accuracy”. Na końcu wystarczy nacisnąć przycisk „Generate”.

### 11.2.2.4 Generacja raportu na życzenie {#generacja-raportu-na-życzenie}

Wymagania funkcjonalne: REQ08

Proces generacji raportu jest analogiczny do opisywanego procesu generacji raportu na żądanie. W momencie widoku konfiguracji raportu należy zmienić wartość w sekcji “Generation type” na “ON_DEMAND”:

<figure>
    <img src="/user-interface/report-config-on-demand.png">
    <figcaption> Widok konfiguracji raportu dla raportów na życzenie [źródło opracowanie własne]</figcaption>
</figure>

Po zmianie wartości wspomnianego pola pojawi się sekcja „Data Range”, w której użytkownik może określić okres, z którego mają pochodzić logi wykorzystane do generacji raportu. Pozostała konfiguracja nie różni się od konfiguracji raportów cyklicznych. Po zakończeniu wszystkich ustawień wystarczy kliknąć przycisk „Generate”.

### 11.2.2.5 Konfiguracja kanałów powiadomień {#konfiguracja-kanałów-powiadomień}

Wymagania funkcjonalne: REQ13, REQ14

Konfiguracja kanałów powiadomień jest dostępna z poziomu widoku „Notifications”. Aby przejść do tego widoku, użytkownik musi zalogować się do systemu, a następnie wybrać zakładkę „Notifications”.

<figure>
    <img src="/user-interface/notification-page.png">
    <figcaption> Widok konfiguracji powiadomień [źródło opracowanie własne]</figcaption>
</figure>

Załóżmy, że użytkownik chce dodać nowy adres e-mail, na który będzie otrzymywał powiadomienia z systemu. W tym celu powinien kliknąć ikonę „plusa” znajdującą się w prawym górnym rogu sekcji „Email”.

<figure>
    <img src="/user-interface/window-for-new-email-address.png">
    <figcaption> Okno służące dodaniu nowego adresu e-mail </figcaption>
</figure>

Po kliknięciu ikony na ekranie wyświetli się okno dialogowe, w którym użytkownik powinien wprowadzić nazwę oraz adres e-mail. Po uzupełnieniu tych danych należy zatwierdzić dodanie adresu e-mail, klikając przycisk „Submit”.

<figure>
    <img src="/user-interface/email-section.png">
    <figcaption> Sekcja e-mail [źródło opracowanie własne]</figcaption>
</figure>

Po zatwierdzeniu użytkownik powinien zobaczyć nowo dodany adres e-mail na liście kanałów powiadomień. Aby przetestować działanie kanału, użytkownik może kliknąć ikonę „Wysyłania wiadomości” w odpowiednim wierszu. Po wykonaniu tej operacji na ekranie powinna pojawić się informacja o pomyślnym wysłaniu wiadomości, a na skrzynkę elektroniczną użytkownika powinien zostać dostarczony e-mail testowy.

<figure>
    <img src="/user-interface/email-section-with-message.png">
    <figcaption> Sekcja “Email” z wiadomością o pomyślnym wysłaniu wiadomości [źródło opracowanie własne]</figcaption>
</figure>

<figure>
    <img src="/user-interface/test-email.png">
    <figcaption> Otrzymany e-mail testowy [źródło opracowanie własne]</figcaption>
</figure>

Jeśli użytkownik chce zmodyfikować dane kanału, powinien kliknąć ikonę „Ołówek”. Po jej wybraniu na ekranie pojawi się okno umożliwiające wprowadzenie nowych danych.

<figure>
    <img src="/user-interface/edit-email.png">
    <figcaption> Okno do edycji e-maila [źródło opracowanie własne]</figcaption>
</figure>

Aby zatwierdzić wprowadzone zmiany, należy kliknąć przycisk „Save”. W przypadku pomyślnej modyfikacji adresu e-mail na ekranie pojawi się komunikat potwierdzający zmianę, analogiczny do tego wyświetlanego podczas wysyłania testowej wiadomości. Dodatkowo zmodyfikowane dane będą widoczne w tabeli.
Możliwość usunięcia kanału powiadomień dostępna jest poprzez kliknięcie ikony „Kosz” w wierszu danego kanału. Po pomyślnym usunięciu kanału na ekranie pojawi się odpowiedni komunikat potwierdzający operację.

<figure>
    <img src="/user-interface/delete-email-message.png">
    <figcaption> Komunikat o usunięciu kanału powiadomień [źródło opracowanie własne]</figcaption>
</figure>

# 12. Bibliografia {#bibliografia}

### Bibliografia

1. <a id="ref1"></a>[The Kubernetes Authors, Home](https://kubernetes.io), dostęp 24 listopada 2024.
2. <a id="ref2"></a>[Docker, Home](https://www.docker.com), dostęp 24 listopada 2024.
3. <a id="ref3"></a>[The Go Authors, Home](https://go.dev), dostęp 24 listopada 2024.
4. <a id="ref4"></a>[Uber Technologies, Fx](https://github.com/uber-go/fx), dostęp 24 listopada 2024.
5. <a id="ref5"></a>[Oracle, Java](https://www.oracle.com/java/), dostęp 24 listopada 2024.
6. <a id="ref6"></a>[Pivotal Software, Spring Boot](https://spring.io/projects/spring-boot), dostęp 24 listopada 2024.
7. <a id="ref7"></a>[Microsoft, TypeScript](https://www.typescriptlang.org), dostęp 24 listopada 2024.
8. <a id="ref8"></a>[Meta, React](https://react.dev), dostęp 24 listopada 2024.
9. <a id="ref9"></a>[Hampton Catlin, Sass](https://sass-lang.com), dostęp 24 listopada 2024.
10. <a id="ref10"></a>[Evan You, Vite](https://vitejs.dev), dostęp 24 listopada 2024.
11. <a id="ref11"></a>[PostgreSQL Global Development Group, PostgreSQL](https://www.postgresql.org), dostęp 24 listopada 2024.
12. <a id="ref12"></a>[MongoDB Inc., MongoDB](https://www.mongodb.com), dostęp 24 listopada 2024.
13. <a id="ref13"></a>[Apache Software Foundation, Kafka](https://kafka.apache.org), dostęp 24 listopada 2024.
14. <a id="ref14"></a>[Elastic N.V., ElasticSearch](https://www.elastic.co/elasticsearch), dostęp 24 listopada 2024.
15. <a id="ref15"></a>[Redis Ltd., Redis](https://redis.io), dostęp 24 listopada 2024.
16. <a id="ref16"></a>[Nginx, Inc., NGINX](https://nginx.org), dostęp 24 listopada 2024.
17. <a id="ref17"></a>[Dynatrace](https://www.dynatrace.com), dostęp 7 grudnia 2024.
18. <a id="ref18"></a>[Datadog](https://www.datadoghq.com), dostęp 7 grudnia 2024.
19. <a id="ref19"></a>[Logz.io](https://logz.io), dostęp 7 grudnia 2024.
20. <a id="ref20"></a>[Azure Monitor](https://azure.microsoft.com/pl-pl/products/monitor), dostęp 7 grudnia 2024.
21. <a id="ref21"></a>[Amazon CloudWatch](https://aws.amazon.com/cloudwatch), dostęp 7 grudnia 2024.
22. <a id="ref22"></a>[Science Logic](https://sciencelogic.com), dostęp 7 grudnia 2024.
23. <a id="ref23"></a>[Elasticsearch](https://www.elastic.co/elasticsearch), dostęp 7 grudnia 2024.
24. <a id="ref24"></a>[Elastic Cloud](https://www.elastic.co/cloud), dostęp 7 grudnia 2024.
25. <a id="ref25"></a>[Gain insights into Kubernetes errors with Elastic Observability logs and OpenAI](https://www.elastic.co/observability-labs/blog/kubernetes-errors-observability-logs-openai), dostęp 9 grudnia 2024.
26. <a id="ref26"></a>[LogEval: A Comprehensive Benchmark Suite for Large Language Models In Log Analysis](https://arxiv.org/pdf/2407.01896).
27. <a id="ref27"></a>[Log Analyzer](https://chatgpt.com/g/g-JBmeWsY2w-log-analyzer), dostęp 7 grudnia 2024.
