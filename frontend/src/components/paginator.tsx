import { Dispatch, SetStateAction } from "react";
import { Pagination } from "react-bootstrap";

export default function Paginator(
    { page, setPage, maxPage }:
    { page: number, setPage: Dispatch<SetStateAction<number>>, maxPage: number }
) {
    return (
        <Pagination className="mt-3 mx-auto max-w-40">
            <Pagination.First onClick={() => setPage(1)} disabled={maxPage == 0} />
            <Pagination.Prev onClick={() => setPage((prev) => prev - 1)} disabled={page == 1} />
            <Pagination.Ellipsis hidden={page < 3} />

            <Pagination.Item hidden={maxPage <= 1} active={page == 1}>{page == 1 ? 1 : page - 1}</Pagination.Item>
            <Pagination.Item active={(page != 1 && page != maxPage) || (page == 1 && maxPage < 3)}>{page}</Pagination.Item>
            <Pagination.Item hidden={maxPage < 3} active={page == maxPage}>{page == maxPage ? maxPage : page + 1}</Pagination.Item>

            <Pagination.Ellipsis hidden={page > maxPage - 3} />
            <Pagination.Next onClick={() => setPage((prev) => prev + 1)} disabled={page == maxPage || maxPage == 0} />
            <Pagination.Last onClick={() => setPage(maxPage)} disabled={maxPage == 0} />
        </Pagination>
    )
}