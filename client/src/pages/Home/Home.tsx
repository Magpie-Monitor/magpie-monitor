import ImportantFindings from './components/ImportantFindings';
import PlaceholderComponent from '@/components/PlaceholderComponent/PlaceholderComponent.tsx';
import SettingsIcon from '@/assets/settings_icon.svg';
import './Home.scss';
import PlaceholderComponentTitle
    from '@/components/PlaceholderComponent/PlaceholderComponentTitle/PlaceholderComponentTitle.tsx';
import ScanStats from './components/ScanStats';

const Home = () => {
    return (
        <div className="home">
            <div className="content">
                <p className="heading">Dashboard</p>
                <div className="dashboard">
                    <PlaceholderComponent
                        icon={<img src={SettingsIcon} alt="Settings Icon"/>}
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
    );
};

export default Home;
