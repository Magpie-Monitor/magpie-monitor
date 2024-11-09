import React, { useState, useEffect } from 'react';
import Table, { TableColumn, TableRow } from 'components/Table/Table';
import Checkbox from 'components/Checkbox/Checkbox';
import ActionButton, { ActionButtonColor } from 'components/ActionButton/ActionButton';
import './EntriesSelector.scss';

interface EntriesSelectorProps<T extends TableRow> {
    selectedItems: T[];
    setSelectedItems: React.Dispatch<React.SetStateAction<T[]>>;
    itemsToExclude: T[];
    onAdd: () => void;
    onClose: () => void;
    fetchData: () => Promise<T[]>;
    columns: TableColumn<T & { key: string }>[];
    getUniqueKey: (item: T) => string;
    entityLabel: string;
    noEntriesMessage?: React.ReactNode;
}

const EntriesSelector = <T extends TableRow>({
                                                 selectedItems,
                                                 setSelectedItems,
                                                 itemsToExclude,
                                                 onAdd,
                                                 onClose,
                                                 fetchData,
                                                 columns,
                                                 getUniqueKey,
                                                 entityLabel,
                                                 noEntriesMessage,
                                             }: EntriesSelectorProps<T>): JSX.Element => {
    const [items, setItems] = useState<T[]>([]);
    const [selectAll, setSelectAll] = useState(false);

    useEffect(() => {
        const fetchItems = async () => {
            try {
                const data = await fetchData();
                setItems(data);
            } catch (error) {
                console.error(`Failed to fetch ${entityLabel}s:`, error);
            }
        };
        fetchItems();
    }, [fetchData, entityLabel]);

    const availableItems = items.filter(
        (item) => !itemsToExclude.some((excluded) => getUniqueKey(excluded) === getUniqueKey(item))
    );

    useEffect(() => {
        setSelectAll(
            availableItems.length > 0 && selectedItems.length === availableItems.length
        );
    }, [selectedItems, availableItems]);

    const handleSelectAllChange = () => {
        setSelectedItems(selectAll ? [] : availableItems);
        setSelectAll(!selectAll);
    };

    const handleCheckboxChange = (item: T) => {
        setSelectedItems((prevSelected) => {
            const isSelected = prevSelected.some(
                (selectedItem) => getUniqueKey(selectedItem) === getUniqueKey(item)
            );
            return isSelected
                // eslint-disable-next-line max-len
                ? prevSelected.filter((selectedItem) => getUniqueKey(selectedItem) !== getUniqueKey(item))
                : [...prevSelected, item];
        });
    };

    const updatedColumns: TableColumn<T & { key: string }>[] = [
        {
            header: (
                <Checkbox
                    checked={selectAll}
                    onChange={handleSelectAllChange}
                />
            ),
            columnKey: 'checkbox',
            customComponent: (row: T & { key: string }) => (
                <Checkbox
                    checked={selectedItems.some(
                        (selectedItem) => getUniqueKey(selectedItem) === getUniqueKey(row)
                    )}
                    onChange={() => handleCheckboxChange(row)}
                />
            ),
        },
        ...columns,
    ];

    return (
        <div className="entries-selector">
            {availableItems.length === 0 ? (
                <div className="entries-selector__no-entries-message">
                    {noEntriesMessage || <p>No {entityLabel} to add.</p>}
                </div>
            ) : (
                <Table
                    columns={updatedColumns}
                    rows={availableItems.map((item) => ({
                        ...item,
                        key: getUniqueKey(item),
                    }))}
                    maxHeight="65vh"
                />
            )}
            <div className="entries-selector__buttons">
                {availableItems.length > 0 && (
                    <ActionButton
                        onClick={onAdd}
                        description="Add"
                        color={ActionButtonColor.GREEN}
                    />
                )}
                <ActionButton
                    onClick={onClose}
                    description="Close"
                    color={ActionButtonColor.RED}
                />
            </div>
        </div>
    );
};

export default EntriesSelector;
