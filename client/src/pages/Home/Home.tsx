import ImportantFindings from './components/ImportantFindings/ImportantFindings';
import SectionComponent from 'components/SectionComponent/SectionComponent.tsx';
import './Home.scss';
import LastReportTitle from './components/LastReportTitle/LestReportTitle';
import ScanStats from './components/ScanStats/ScanStats';
import PageTemplate from 'components/PageTemplate/PageTemplate';
import HeaderWithIcon from 'components/PageTemplate/components/HeaderWithIcon/HeaderWithIcon';
import SVGIcon from 'components/SVGIcon/SVGIcon';

const Home = () => {
  return (
    <PageTemplate header={<HeaderWithIcon title={'Dashboard'} />}>
      <SectionComponent
        icon={<SVGIcon iconName="chart-icon" />}
        title={
          <LastReportTitle
            source="production-services"
            startTime="19.04.2023"
            endTime="25.04.2023"
          />
        }
      >
        <div>
          <ScanStats />
          <ImportantFindings />
        </div>
      </SectionComponent>
    </PageTemplate>
  );
};

export default Home;
