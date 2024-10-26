import ImportantFindings from './components/ImportantFindings/ImportantFindings';
import SectionComponent from 'components/SectionComponent/SectionComponent.tsx';
import './Home.scss';
import LastReportTitle from './components/LastReportTitle/LestReportTitle';
import ScanStats from './components/ScanStats/ScanStats';

const Home = () => {
  return (
    <div className="home">
      <div className="home__content">
        <div>
          <p className="home__content__heading">Dashboard</p>
          <div>
            <SectionComponent
              icon={'setting-icon'}
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
          </div>
        </div>
      </div>
    </div>
  );
};

export default Home;
