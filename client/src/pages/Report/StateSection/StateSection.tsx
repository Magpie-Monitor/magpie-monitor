import SectionComponent from 'components/SectionComponent/SectionComponent.tsx';
import TagButton from 'components/TagButton/TagButton.tsx';
import { useState } from 'react';
import SVGIcon from 'components/SVGIcon/SVGIcon.tsx';

const StateSection = () => {
    const [state, setState] = useState('enabled');
    return (
        <SectionComponent icon={<SVGIcon iconName='precision-icon' />} title={'State'}>
            <div className="precision-section__input-group">
                <TagButton
                    listItems={['enabled', 'disabled']}
                    chosenItem={state}
                    onSelect={setState}
                />
            </div>
        </SectionComponent>
    );
};
export default StateSection;