import {
  Diagram,
  DiagramType,
  FlowchartDiagram,
  FlowchartNode,
  FlowchartEdge,
  SequenceDiagram,
  SequenceParticipant,
  SequenceMessage,
} from './types';

export function detectDiagramType(input: string): DiagramType {
  const trimmed = input.trim().toLowerCase();
  if (trimmed.startsWith('graph') || trimmed.startsWith('flowchart')) {
    return 'flowchart';
  }
  if (trimmed.startsWith('sequencediagram')) {
    return 'sequence';
  }
  return 'unknown';
}

export function parse(input: string): Diagram {
  const type = detectDiagramType(input);

  switch (type) {
    case 'flowchart':
      return parseFlowchart(input);
    case 'sequence':
      return parseSequence(input);
    default:
      return { type: 'unknown' };
  }
}

function parseFlowchart(input: string): FlowchartDiagram {
  // Split by newlines or semicolons
  const lines = input
    .split(/[\n;]/)
    .map(l => l.trim())
    .filter(l => l);
  const nodes: Map<string, FlowchartNode> = new Map();
  const edges: FlowchartEdge[] = [];

  // Parse direction from first line
  let direction: FlowchartDiagram['direction'] = 'TD';
  const firstLine = lines[0].toLowerCase();
  if (firstLine.includes(' lr')) direction = 'LR';
  else if (firstLine.includes(' rl')) direction = 'RL';
  else if (firstLine.includes(' bt')) direction = 'BT';
  else if (firstLine.includes(' tb') || firstLine.includes(' td')) direction = 'TD';

  // Parse nodes and edges
  for (let i = 1; i < lines.length; i++) {
    const line = lines[i];

    // Skip subgraph, end, style, class directives
    if (/^(subgraph|end|style|class|click|linkstyle)/i.test(line)) continue;

    // Parse chained edges: A --> B --> C or A[Label] --> B[Label] --> C[Label]
    const nodePattern = /(\w+)(\[[^\]]+\]|\([^)]+\)|\{[^}]+\}|\(\([^)]+\)\)|\[\[[^\]]+\]\])?/g;
    const edgePattern = /(-->|---|-\.->|-.->|==>|-->\|[^|]+\||---\|[^|]+\||-\.->\|[^|]+\|)/g;

    // Find all nodes and edges in the line
    const nodeMatches: Array<{ id: string; shape?: string; index: number }> = [];
    const edgeMatches: Array<{ style: string; index: number }> = [];

    let nodeMatch;
    while ((nodeMatch = nodePattern.exec(line)) !== null) {
      // Skip if this is part of an edge pattern
      if (nodeMatch[0].startsWith('-') || nodeMatch[0].startsWith('=')) continue;
      nodeMatches.push({
        id: nodeMatch[1],
        shape: nodeMatch[2],
        index: nodeMatch.index,
      });
    }

    let edgeMatch;
    while ((edgeMatch = edgePattern.exec(line)) !== null) {
      edgeMatches.push({
        style: edgeMatch[1],
        index: edgeMatch.index,
      });
    }

    // If we found edges, process them
    if (edgeMatches.length > 0 && nodeMatches.length >= 2) {
      // Sort by index position
      nodeMatches.sort((a, b) => a.index - b.index);
      edgeMatches.sort((a, b) => a.index - b.index);

      // Add all nodes
      for (const nm of nodeMatches) {
        if (!nodes.has(nm.id)) {
          nodes.set(nm.id, parseNodeShape(nm.id, nm.shape));
        } else if (nm.shape && nodes.get(nm.id)!.label === nm.id) {
          nodes.set(nm.id, parseNodeShape(nm.id, nm.shape));
        }
      }

      // Create edges between consecutive nodes
      for (let j = 0; j < edgeMatches.length && j < nodeMatches.length - 1; j++) {
        const fromNode = nodeMatches[j];
        const toNode = nodeMatches[j + 1];
        const edgeStyle = edgeMatches[j].style;

        // Parse edge label
        let edgeLabel: string | undefined;
        const labelMatch = edgeStyle.match(/\|([^|]+)\|/);
        if (labelMatch) {
          edgeLabel = labelMatch[1];
        }

        // Determine edge style
        let style: FlowchartEdge['style'] = 'solid';
        let arrow: FlowchartEdge['arrow'] = 'arrow';

        if (edgeStyle.includes('-.') || edgeStyle.includes('-.-')) {
          style = 'dotted';
        } else if (edgeStyle.includes('==')) {
          style = 'thick';
        }

        if (edgeStyle.includes('---')) {
          arrow = 'none';
        }

        edges.push({ from: fromNode.id, to: toNode.id, label: edgeLabel, style, arrow });
      }
      continue;
    }

    // Parse standalone node definition: A[Label] or A(Label) etc.
    const standaloneMatch = line.match(/^(\w+)(\[[^\]]+\]|\([^)]+\)|\{[^}]+\}|\(\([^)]+\)\)|\[\[[^\]]+\]\])$/);
    if (standaloneMatch) {
      const nodeId = standaloneMatch[1];
      const nodeShape = standaloneMatch[2];
      nodes.set(nodeId, parseNodeShape(nodeId, nodeShape));
    }
  }

  return {
    type: 'flowchart',
    direction,
    nodes: Array.from(nodes.values()),
    edges,
  };
}

function parseNodeShape(id: string, shapeStr?: string): FlowchartNode {
  if (!shapeStr) {
    return { id, label: id, shape: 'rectangle' };
  }

  let label = id;
  let shape: FlowchartNode['shape'] = 'rectangle';

  if (shapeStr.startsWith('((') && shapeStr.endsWith('))')) {
    // ((Circle))
    label = shapeStr.slice(2, -2);
    shape = 'circle';
  } else if (shapeStr.startsWith('([') && shapeStr.endsWith('])')) {
    // ([Stadium])
    label = shapeStr.slice(2, -2);
    shape = 'stadium';
  } else if (shapeStr.startsWith('[[') && shapeStr.endsWith(']]')) {
    // [[Subroutine]]
    label = shapeStr.slice(2, -2);
    shape = 'rectangle';
  } else if (shapeStr.startsWith('(') && shapeStr.endsWith(')')) {
    // (Rounded)
    label = shapeStr.slice(1, -1);
    shape = 'rounded';
  } else if (shapeStr.startsWith('{') && shapeStr.endsWith('}')) {
    // {Diamond}
    label = shapeStr.slice(1, -1);
    shape = 'diamond';
  } else if (shapeStr.startsWith('[') && shapeStr.endsWith(']')) {
    // [Rectangle]
    label = shapeStr.slice(1, -1);
    shape = 'rectangle';
  }

  return { id, label, shape };
}

function parseSequence(input: string): SequenceDiagram {
  // Split by newlines or semicolons
  const lines = input
    .split(/[\n;]/)
    .map(l => l.trim())
    .filter(l => l);
  const participants: Map<string, SequenceParticipant> = new Map();
  const messages: SequenceMessage[] = [];

  for (let i = 1; i < lines.length; i++) {
    const line = lines[i];

    // Skip notes, activations, loops, etc.
    if (/^(note|activate|deactivate|loop|end|alt|else|opt|par|and|rect|critical|break)/i.test(line)) {
      continue;
    }

    // Parse participant definition
    const participantMatch = line.match(/^participant\s+(\w+)(?:\s+as\s+(.+))?$/i);
    if (participantMatch) {
      const id = participantMatch[1];
      const label = participantMatch[2] || id;
      participants.set(id, { id, label });
      continue;
    }

    // Parse actor definition
    const actorMatch = line.match(/^actor\s+(\w+)(?:\s+as\s+(.+))?$/i);
    if (actorMatch) {
      const id = actorMatch[1];
      const label = actorMatch[2] || id;
      participants.set(id, { id, label });
      continue;
    }

    // Parse message: A->>B: Message or A-->>B: Message, A->B:, A-->B:
    const msgMatch = line.match(/^(\w+)\s*(--?>?>|--?>>?|-\.->?)\s*(\w+)\s*:\s*(.*)$/);
    if (msgMatch) {
      const from = msgMatch[1];
      const arrowStyle = msgMatch[2];
      const to = msgMatch[3];
      const label = msgMatch[4];

      // Auto-add participants
      if (!participants.has(from)) {
        participants.set(from, { id: from, label: from });
      }
      if (!participants.has(to)) {
        participants.set(to, { id: to, label: to });
      }

      const style: SequenceMessage['style'] = arrowStyle.includes('--') ? 'dotted' : 'solid';
      const arrow: SequenceMessage['arrow'] = arrowStyle.includes('>>') ? 'arrow' : 'open';

      messages.push({ from, to, label, style, arrow });
    }
  }

  return {
    type: 'sequence',
    participants: Array.from(participants.values()),
    messages,
  };
}
