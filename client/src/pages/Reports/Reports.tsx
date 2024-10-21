import './Reports.scss';
import PlaceholderComponent from '@/components/PlaceholderComponent/PlaceholderComponent.tsx';
// import Table from '@/components/Table/Table.tsx';

const Reports = () => {
    return (
        <div className="reports">
            <div className="reports__content">
                <p className="reports__content__heading">Reports</p>
                <div className="reports__content__dashboard">
                    <PlaceholderComponent
                        icon={'setting-icon'}
                        title={<p> Weekly reports</p>}>
                        <div>
                            {/*<Table></Table>*/}
                        </div>
                    </PlaceholderComponent>
                </div>
            </div>
        </div>
    );
};

export default Reports;
