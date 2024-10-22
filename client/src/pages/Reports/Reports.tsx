import './Reports.scss';
import SectionComponent from 'components/SectionComponent/SectionComponent.tsx';
// import Table from '@/components/Table/Table.tsx';

const Reports = () => {
    return (
        <div className="reports">
            <div className="reports__content">
                <p className="reports__content__heading">Reports</p>
                <div className="reports__content__dashboard">
                    <SectionComponent
                        icon={'setting-icon'}
                        title={<p> Weekly reports</p>}
                    >
                        <div>{/*<Table></Table>*/}</div>
                    </SectionComponent>
                </div>
            </div>
        </div>
    );
};

export default Reports;
