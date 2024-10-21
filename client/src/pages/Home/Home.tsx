import ImportantFindings from './components/ImportantFindings';
import PlaceholderComponent from '@/components/PlaceholderComponent/PlaceholderComponent.tsx';
import './Home.scss';
import PlaceholderComponentTitle
    from '@/components/PlaceholderComponent/PlaceholderComponentTitle/PlaceholderComponentTitle.tsx';
import ScanStats from './components/ScanStats';

//panel name should be aligned to the left of the content

const Home = () => {
    return (
        <div className="home">
            <div className="home__content">
                <div>
                    <p className="home__content__heading">Dashboard</p>
                    <div>
                        <PlaceholderComponent
                            icon={'setting-icon'}
                            title={<PlaceholderComponentTitle source="production-services"
                                                              startTime="19.04.2023" endTime="25.04.2023"/>}>
                            <div>
                                <ScanStats/>
                                <ImportantFindings/>
                            </div>
                        </PlaceholderComponent>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Home;
