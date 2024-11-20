import SectionComponent from 'components/SectionComponent/SectionComponent.tsx';
import TagButton from 'components/TagButton/TagButton.tsx';
import React, { useState } from 'react';
import SVGIcon from 'components/SVGIcon/SVGIcon.tsx';
import {ReportType} from 'api/managment-service.ts';

interface ReportGenerationTypeProps {
    setParentGenerationType: (generationType: ReportType) => void;
}

const ReportGenerationType: React.FC<ReportGenerationTypeProps> = ({ setParentGenerationType }) => {
    const [generationType, setGenerationType] = useState<ReportType>('ON-DEMAND');

    const handleGenerationTypeChange = (newGenerationType: ReportType) => {
        setGenerationType(newGenerationType);
        setParentGenerationType(newGenerationType);
    };

    return (
        <SectionComponent icon={<SVGIcon iconName='precision-icon' />} title={'Generation type'}>
            <div className="precision-section__input-group">
                <TagButton
                    listItems={['ON-DEMAND', 'SCHEDULED']}
                    chosenItem={generationType}
                    onSelect={handleGenerationTypeChange}
                />
            </div>
        </SectionComponent>
    );
};
export default ReportGenerationType;