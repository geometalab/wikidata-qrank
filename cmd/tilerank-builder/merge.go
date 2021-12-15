// SPDX-License-Identifier: MIT

package main

import (
	"bufio"
	"context"
	"io"
	"strconv"
)

// TODO: Implement k-way merge of sorted TileCount files.
func mergeTileCounts(r []io.Reader, out chan<- TileCount, ctx context.Context) error {
	defer close(out)
	if len(r) == 0 {
		return nil
	}

	if len(r) > 1 {
		// Add readers to heap
		tch := make(TileCountHeap, 0)
		for i := 0; i < len(r); i++ {
			reader := r[i]
			scanner := bufio.NewScanner(reader)
			if scanner.Scan() {
				firstLine := scanner.Text()
				tch.Push(&TileCountReader{scanner: *scanner, tc: TileCountFromString(firstLine)})
			}
		}

		// Do until heap is empty
		for tch.Len() != 0 {
			toPush := tch[0].tc
			print(toPush.ToString(), " \n")
			out <- toPush
			// If reader still has lines
			if tch[0].scanner.Scan() {
				nextLine := tch[0].scanner.Text()
				tch.fix(tch[0], TileCountFromString(nextLine), tch[0].scanner)
			} else {
				tch.Pop()
			}
		}
	}
	print("\n")

	scanner := bufio.NewScanner(r[len(r)-1])
	for scanner.Scan() {
		// Check if our task has been canceled. Typically, this can happen
		// because of an error in another goroutine in the same x.sync.errroup.
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		match := tileLogRegexp.FindStringSubmatch(scanner.Text())
		if match == nil || len(match) != 5 {
			continue
		}
		zoom, _ := strconv.Atoi(match[1])
		if zoom < 0 {
			continue
		}
		x, _ := strconv.ParseUint(match[2], 10, 32)
		y, _ := strconv.ParseUint(match[3], 10, 32)
		count, _ := strconv.ParseUint(match[4], 10, 64)
		key := MakeTileKey(uint8(zoom), uint32(x), uint32(y))
		out <- TileCount{Key: key, Count: count}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
