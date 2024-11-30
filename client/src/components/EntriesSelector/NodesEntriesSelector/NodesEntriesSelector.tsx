import React from 'react';
import EntriesSelector from 'components/EntriesSelector/EntriesSelector';
import LinkComponent from 'components/LinkComponent/LinkComponent.tsx';
import {NodeDataRow} from 'pages/Report/NodesSection/NodesSection.tsx';
import KindTag from 'components/KindTag/KindTag.tsx';

interface NodesEntriesSelectorProps {
  selectedNodes: NodeDataRow[];
  setSelectedNodes: React.Dispatch<React.SetStateAction<NodeDataRow[]>>;
  nodesToExclude: NodeDataRow[];
  onAdd: () => void;
  onClose: () => void;
  availableNodes: NodeDataRow[];
}

const NodesEntriesSelector: React.FC<NodesEntriesSelectorProps> = ({
                                                                     selectedNodes,
                                                                     setSelectedNodes,
                                                                     nodesToExclude,
                                                                     onAdd,
                                                                     onClose,
                                                                     availableNodes
                                                                   }) => {

  const getUniqueKey = (node: NodeDataRow) => node.name;

  const columns = [
    {
      header: 'Name',
      columnKey: 'name',
      customComponent: (row: NodeDataRow) => (
        <LinkComponent isRunning={row.running}>
          {row.name}
        </LinkComponent>
      ),
    },
    {
      header: 'Kind',
      columnKey: '',
      customComponent: () => (
        <KindTag
          name={'Node'}
        />
      ),
    },
  ];

  return (
    <EntriesSelector<NodeDataRow>
      selectedItems={selectedNodes}
      setSelectedItems={setSelectedNodes}
      itemsToExclude={nodesToExclude}
      onAdd={onAdd}
      onClose={onClose}
      columns={columns}
      items={availableNodes}
      getKey={getUniqueKey}
      entityLabel="node"
      noEntriesMessage={<p>There is no node to add.</p>}
      title="Select Nodes"
    />
  );
};

export default NodesEntriesSelector;
