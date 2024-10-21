import './Settings.scss';
import PlaceholderComponent from '@/components/PlaceholderComponent/PlaceholderComponent.tsx';
// import Table from '@/components/Table/Table.tsx';

const Settings = () => {
    return (
        <div className="settings">
            <div className="settings__content">
                <p className="settings__content__heading">Notification channels</p>
                <div className="settings__content__dashboard">
                    <PlaceholderComponent
                        icon={'setting-icon'}
                        title={<p> Slack</p>}>
                        <div>
                            {/*<Table></Table>*/}
                        </div>
                    </PlaceholderComponent>
                </div>
                <div className="settings__content__dashboard">
                    <PlaceholderComponent
                        icon={'setting-icon'}
                        title={<p> Discord</p>}>
                        <div>
                            {/*<Table></Table>*/}
                        </div>
                    </PlaceholderComponent>
                </div>
                <div className="settings__content__dashboard">
                    <PlaceholderComponent
                        icon={'setting-icon'}
                        title={<p> Email</p>}>
                        <div>
                            {/*<Table></Table>*/}
                        </div>
                    </PlaceholderComponent>
                </div>
            </div>
        </div>
    );
};

export default Settings;
