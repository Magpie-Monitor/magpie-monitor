import './Settings.scss';
import SectionComponent from '@/components/SectionComponent/SectionComponent.tsx';
// import Table from '@/components/Table/Table.tsx';

const Settings = () => {
    return (
        <div className="settings">
            <div className="settings__content">
                <p className="settings__content__heading">Notification channels</p>
                <div className="settings__content__dashboard">
                    <SectionComponent
                        icon={'setting-icon'}
                        title={<p> Slack</p>}>
                        <div>
                            {/*<Table></Table>*/}
                        </div>
                    </SectionComponent>
                </div>
                <div className="settings__content__dashboard">
                    <SectionComponent
                        icon={'setting-icon'}
                        title={<p> Discord</p>}>
                        <div>
                            {/*<Table></Table>*/}
                        </div>
                    </SectionComponent>
                </div>
                <div className="settings__content__dashboard">
                    <SectionComponent
                        icon={'setting-icon'}
                        title={<p> Email</p>}>
                        <div>
                            {/*<Table></Table>*/}
                        </div>
                    </SectionComponent>
                </div>
            </div>
        </div>
    );
};

export default Settings;
